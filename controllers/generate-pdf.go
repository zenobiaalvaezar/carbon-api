package controllers

import (
	"bytes"
	"carbon-api/models"
	"carbon-api/repositories"
	"carbon-api/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type GeneratePdfController struct {
	UserRepository           repositories.UserRepository
	CarbonElectricRepository repositories.CarbonElectricRepository
	CarbonSummaryRepository  repositories.CarbonSummaryRepository
}

func NewGeneratePdfController(userRepository repositories.UserRepository, carbonElectricRepository repositories.CarbonElectricRepository, carbonSummaryRepository repositories.CarbonSummaryRepository) *GeneratePdfController {
	return &GeneratePdfController{UserRepository: userRepository, CarbonElectricRepository: carbonElectricRepository, CarbonSummaryRepository: carbonSummaryRepository}
}

type ReportData struct {
	Address      string
	ReportDate   string
	Emission     string
	FuelData     []Fuel
	EmissionData EmissionData
}

type Fuel struct {
	ID             int
	Category       string
	Name           string
	EmissionFactor float64
	Price          float64
	Unit           string
	Value          float64
	X              float64
	Y              float64
	TextX, TextY   float64
}

type EmissionData struct {
	NationalAvg int
	ProvinceAvg int
}

type ReportDataSummary struct {
	TotalEmission int
	UserName      string
	UserEmail     string
	TreesNeeded   int
	LastRecords   []EmissionRecord
	AIPredictions string
}

type EmissionRecord struct {
	FuelEmission     int
	ElectricEmission int
	TotalEmission    int
	TotalTree        int
}

func generateContent(totalEmission float64, treesNeeded int, lastRecords []EmissionRecord) string {
	ctx := context.Background()

	// Set up Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// Constructing the dynamic prompt
	var lastRecordsStr string
	for i, record := range lastRecords {
		lastRecordsStr += fmt.Sprintf("Rekor %d: Emisi Bahan Bakar = %.2f, Emisi Listrik = %.2f, Total Emisi = %.2f\n",
			i+1, record.FuelEmission, record.ElectricEmission, record.TotalEmission)
	}

	prompt := fmt.Sprintf(`
		Analisis data emisi karbon berikut dan berikan prediksi untuk emisi di masa depan serta rekomendasi untuk mencapai netralitas karbon. Prediksi harus didasarkan pada catatan sebelumnya dan mencakup langkah-langkah yang dapat dilakukan.
		Data:
		
		Total Emisi: %.2f
		Jumlah Pohon yang Dibutuhkan: %d
		Catatan Terakhir:
		%s
		
		Format Output:
		Prediksi emisi untuk 6 bulan ke depan.
		Garis waktu untuk mencapai netralitas karbon jika tren saat ini berlanjut.
		Rekomendasi untuk mengurangi emisi.
	`, totalEmission, treesNeeded, lastRecordsStr)

	// Request content generation from the model
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	var recommendations string
	reHTML := regexp.MustCompile(`(?i)<!DOCTYPE html>|<(html|head|body)[^>]*>|</(html|head|body)>|<title>.*</title>`)
	reBackticks := regexp.MustCompile("```")

	// Process the response
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					var validJSON interface{}
					// If it's valid JSON, try to unmarshal and return as list of recommendations
					if err := json.Unmarshal([]byte(txt), &validJSON); err == nil {
						var recipes []string
						if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
							log.Fatal(err)
						}
						recommendations += strings.Join(recipes, "\n") + "\n"
					} else {
						// Clean HTML tags and backticks if not valid JSON
						plainText := reHTML.ReplaceAllString(string(txt), "")
						plainText = reBackticks.ReplaceAllString(plainText, "")
						recommendations += plainText + "\n"
					}
				}
			}
		}
	}

	return recommendations
}

// PdfHandlerSummary godoc
// @Summary Generate a summary PDF report
// @Description Generate a summary PDF report with emission data and AI-generated predictions, then send it via email
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {string} string "PDF generated and sent to fr081938@gmail.com"
// @Failure 500 {object} map[string]string "Error generating PDF or sending email"
// @Security BearerAuth
// @Router /generate-pdf-summary [post]
func (ctrl *GeneratePdfController) PdfHandlerSummary(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User ID is invalid or missing. Please log in to continue.",
		})
	}

	user, _, err := ctrl.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	lastRecords, _, err := ctrl.CarbonElectricRepository.GetLast3CarbonElectrics(userID)
	if err != nil {
		return err
	}

	// Attempt to get the carbon summary
	carbonSummary, _, err := ctrl.CarbonSummaryRepository.GetCarbonSummary(userID)

	// Check if the error is due to "Carbon summary not found"
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "Carbon summary not found" {
			// Log the error but continue with a fallback or default data
			log.Printf("Warning: Carbon summary not found for user %d. Using default values.", userID)
			// You can initialize carbonSummary with default values here if needed
			carbonSummary = models.CarbonSummaryResponse{
				TotalEmission: 0,
				TotalTree:     0,
			}
		} else {
			// If it's another error, return it
			return err
		}
	}

	tmplPath := filepath.Join("templates", "weekly_email_emission_record_pdf.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to load template")
	}

	var emissionRecords []EmissionRecord
	for _, record := range lastRecords {
		totalTree := utils.CalculateTotalTree(record.TotalConsumption)

		emissionRecords = append(emissionRecords, EmissionRecord{
			FuelEmission:     int(record.UsageAmount),
			ElectricEmission: int(record.TotalConsumption),
			TotalEmission:    int(record.TotalConsumption),
			TotalTree:        totalTree,
		})
	}

	var renderedHTML bytes.Buffer
	data := ReportDataSummary{
		TotalEmission: int(carbonSummary.TotalEmission),
		UserName:      user.Name,
		UserEmail:     user.Email,
		TreesNeeded:   carbonSummary.TotalTree,
		LastRecords:   emissionRecords,
	}

	aiPredictions := generateContent(carbonSummary.TotalEmission, carbonSummary.TotalTree, emissionRecords)
	data.AIPredictions = aiPredictions

	if err := tmpl.Execute(&renderedHTML, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to render template")
	}

	pdfData, err := utils.HtmlToPDF(renderedHTML.String())
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to generate PDF")
	}

	subject := "Your Generated PDF Report"
	dataBody := map[string]string{
		"Name":          user.Name,
		"TotalEmission": strconv.FormatFloat(carbonSummary.TotalEmission, 'f', -1, 64),
		"Email":         user.Email,
		"TotalTrees":    string(carbonSummary.TotalTree),
	}

	emailBody, err := utils.RenderTemplate(dataBody, "templates/weekly_email_emission_record.html")
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	if err := utils.SendEmailWithPdfAttachment(user.Email, subject, emailBody, pdfData); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to send email")
	}

	return c.String(http.StatusOK, "PDF generated and sent to "+user.Email)
}
