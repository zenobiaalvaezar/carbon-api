package services

import (
	"bytes"
	"carbon-api/repositories"
	"carbon-api/utils"
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type IGeneratePdfService interface {
	PdfHandler(userID int) error
	PdfHandlerSummary(userID int) error
}

type GeneratePdfService struct {
	UserRepository     repositories.UserRepository
	ElectricRepository repositories.ElectricRepository
	FuelRepository     repositories.FuelRepository
}

func NewGeneratePdfService(userRepository repositories.UserRepository, electricRepository repositories.ElectricRepository, fuelRepository repositories.FuelRepository) IGeneratePdfService {
	return &GeneratePdfService{UserRepository: userRepository, ElectricRepository: electricRepository, FuelRepository: fuelRepository}
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
	NationalAvg float64
	ProvinceAvg float64
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

func (ctrl *GeneratePdfService) PdfHandler(userID int) error {
	user, _, err := ctrl.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	nationalAvg, provinceAvg, err := ctrl.ElectricRepository.GetAverageEmission(user.ProvinceID)
	if err != nil {
		return err
	}

	fuels, _, err := ctrl.FuelRepository.GetTop4FuelsByEmissionFactor()
	if err != nil {
		return err
	}

	go func() {
		tmplPath := filepath.Join("templates", "carbon_emission_per_province_pdf.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Printf("Error loading template: %v", err)
		}

		var fuelData []Fuel
		for _, fuel := range fuels {
			fuelData = append(fuelData, Fuel{
				ID:             fuel.ID,
				Category:       fuel.Category,
				Name:           fuel.Name,
				EmissionFactor: fuel.EmissionFactor,
				Price:          fuel.Price,
				Unit:           fuel.Unit,
				Value:          fuel.EmissionFactor,
			})
		}

		for i, fuel := range fuelData {
			fuel.X = float64(i*90 + 40)
			fuel.Y = 200 - fuel.Value
			fuel.TextX = fuel.X + 35
			fuel.TextY = fuel.Y + (fuel.Value / 2)
			fuelData[i] = fuel
		}

		const maxValue = 50.0
		nationalAvg = scaleValueToMax(nationalAvg, maxValue)
		provinceAvg = scaleValueToMax(provinceAvg, maxValue)

		var renderedHTML bytes.Buffer
		data := ReportData{
			Address:    user.Address,
			ReportDate: time.Now().Format("January 2, 2006"),
			Emission:   "CO₂ Emission",
			FuelData:   fuelData,
			EmissionData: EmissionData{
				NationalAvg: nationalAvg,
				ProvinceAvg: provinceAvg,
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
		dataBofy := map[string]string{
			"Name": "Fathur Rohman Wahidd",
		}

		emailBody, err := utils.RenderTemplate(dataBofy, "templates/regist_success.html")
		if err != nil {
			log.Fatalf("Error rendering template: %v", err)
		}

		if err := utils.SendEmailWithPdfAttachment("fr081938@gmail.com", subject, emailBody, pdfData); err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()

	return nil
}

func scaleValueToMax(value, maxValue float64) float64 {
	if value > maxValue {
		// If the value exceeds the maximum, scale it to the max
		return maxValue
	}
	// Otherwise, return the original value
	return value
}

func (ctrl *GeneratePdfService) PdfHandlerSummary(userID int) error {
	tmplPath := filepath.Join("templates", "summary.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
	}

	var renderedHTML bytes.Buffer
	data := ReportDataSummary{
		TotalEmission: 38901,
		UserName:      "Fathur Rohman",
		UserEmail:     "fr081938@gmail.com",
		TreesNeeded:   155,
		LastRecords: []EmissionRecord{
			{FuelEmission: 1905, ElectricEmission: 203, TotalEmission: 109, TotalTree: 120},
			{FuelEmission: 1905, ElectricEmission: 203, TotalEmission: 109, TotalTree: 120},
		},
	}

	aiPredictions := generateContent()
	data.AIPredictions = aiPredictions

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

	return nil
}
