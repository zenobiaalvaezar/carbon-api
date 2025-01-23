package services

import (
	"bytes"
	"carbon-api/repositories"
	"carbon-api/utils"
	"log"
	"path/filepath"
	"text/template"
	"time"
)

type IGeneratePdfService interface {
	PdfHandler(userID int) error
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

		if err := utils.SendEmailWithPdfAttachment(user.Email, subject, emailBody, pdfData); err != nil {
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
