package controllers

import (
	"bytes"
	"carbon-api/utils"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type GeneratePdfController struct{}

func NewGeneratePdfController() *GeneratePdfController {
	return &GeneratePdfController{}
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
	Value          float64 // This will represent the value for the bar height
	X              float64 // X position for the bar
	Y              float64 // Y position for the bar
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
}

func generateContent() string {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(`
		Analisis data emisi karbon berikut dan berikan prediksi untuk emisi di masa depan serta rekomendasi untuk mencapai netralitas karbon. Prediksi harus didasarkan pada catatan sebelumnya dan mencakup langkah-langkah yang dapat dilakukan.
		Data:

		Total Emisi: 38.901
		Jumlah Pohon yang Dibutuhkan: 155
		Catatan Terakhir:
		Rekor 1: Emisi Bahan Bakar = 1.905, Emisi Listrik = 203, Total Emisi = 109
		Rekor 2: Emisi Bahan Bakar = 1.905, Emisi Listrik = 203, Total Emisi = 109
	
		Format Output on single string:

		Prediksi emisi untuk 6 bulan ke depan.
		Garis waktu untuk mencapai netralitas karbon jika tren saat ini berlanjut.
		Rekomendasi untuk mengurangi emisi.
	`,
	))
	if err != nil {
		log.Fatal(err)
	}

	var recommendations string
	reHTML := regexp.MustCompile(`(?i)<!DOCTYPE html>|<(html|head|body)[^>]*>|</(html|head|body)>|<title>.*</title>`)
	reBackticks := regexp.MustCompile("```")

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					var validJSON interface{}
					if err := json.Unmarshal([]byte(txt), &validJSON); err == nil {
						var recipes []string
						if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
							log.Fatal(err)
						}
						recommendations += strings.Join(recipes, "\n") + "\n"
					} else {
						plainText := reHTML.ReplaceAllString(string(txt), "")   // Remove the <html>, <head>, <body> tags
						plainText = reBackticks.ReplaceAllString(plainText, "") // Remove backticks
						recommendations += plainText + "\n"
					}
				}
			}
		}
	}

	return recommendations
}

// PdfHandler godoc
// @Summary Generate a PDF report for carbon emissions
// @Description Generate a PDF report based on emission data and send it via email
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {string} string "PDF generated and sent to fr081938@gmail.com"
// @Failure 500 {object} map[string]string "Error generating PDF or sending email"
// @Security BearerAuth
// @Router /generate-pdf [post]
func (ctrl *GeneratePdfController) PdfHandler(c echo.Context) error {
	go func() {
		tmplPath := filepath.Join("templates", "report.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Printf("Error loading template: %v", err)
		}

		fuelData := []Fuel{
			{ID: 1, Category: "Bahan Bakar Cair", Name: "Pertamax Plus/Turbo", EmissionFactor: 2.368, Price: 13250, Unit: "Liter", Value: 150},
			{ID: 2, Category: "Bahan Bakar Cair", Name: "Pertamax", EmissionFactor: 2.363, Price: 12500, Unit: "Liter", Value: 100},
			{ID: 3, Category: "Bahan Bakar Cair", Name: "Pertalite", EmissionFactor: 2.367, Price: 12000, Unit: "Liter", Value: 80},
			{ID: 4, Category: "Bahan Bakar Cair", Name: "Premium", EmissionFactor: 2.373, Price: 6500, Unit: "Liter", Value: 130},
			{ID: 2, Category: "Bahan Bakar Cair", Name: "Pertamax", EmissionFactor: 2.363, Price: 12500, Unit: "Liter", Value: 100},
			{ID: 3, Category: "Bahan Bakar Cair", Name: "Pertalite", EmissionFactor: 2.367, Price: 12000, Unit: "Liter", Value: 80},
			{ID: 4, Category: "Bahan Bakar Cair", Name: "Premium", EmissionFactor: 2.373, Price: 6500, Unit: "Liter", Value: 130},
		}

		for i, fuel := range fuelData {
			fuel.X = float64(i*90 + 40)
			fuel.Y = 200 - fuel.Value
			fuel.TextX = fuel.X + 35
			fuel.TextY = fuel.Y + (fuel.Value / 2)
			fuelData[i] = fuel
		}

		var renderedHTML bytes.Buffer
		data := ReportData{
			Address:    "123 Greenway Blvd, Ontario",
			ReportDate: time.Now().Format("January 2, 2006"),
			Emission:   "CO₂ Emission",
			FuelData:   fuelData,
			EmissionData: EmissionData{
				NationalAvg: 100,
				ProvinceAvg: 80,
			},
		}

		if err := tmpl.Execute(&renderedHTML, data); err != nil {
			log.Printf("Error rendering template: %v", err)
		}

		pdfData, err := utils.HtmlToPDF(renderedHTML.String())
		if err != nil {
			log.Printf("Error generating PDF: %v", err)
		}

		subject := "Your Generated PDF Report"
		body := "Please find your PDF report attached."
		if err := utils.SendEmailWithPdfAttachment("fr081938@gmail.com", subject, body, pdfData); err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()

	return c.String(http.StatusOK, "PDF generated and sent to fr081938@gmail.com")
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
	tmplPath := filepath.Join("templates", "summary.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to load template")
	}

	var renderedHTML bytes.Buffer
	data := ReportDataSummary{
		TotalEmission: 38901,
		UserName:      "Fathur Rohman",
		UserEmail:     "fr081938@gmail.com",
		TreesNeeded:   155,
		LastRecords: []EmissionRecord{
			{FuelEmission: 1905, ElectricEmission: 203, TotalEmission: 109},
			{FuelEmission: 1905, ElectricEmission: 203, TotalEmission: 109},
		},
	}

	aiPredictions := generateContent()
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
	body := "Please find your PDF report attached."
	if err := utils.SendEmailWithPdfAttachment("fr081938@gmail.com", subject, body, pdfData); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to send email")
	}

	return c.String(http.StatusOK, "PDF generated and sent to fr081938@gmail.com")
}
