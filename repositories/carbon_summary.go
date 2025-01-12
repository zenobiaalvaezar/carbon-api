package repositories

import (
	"carbon-api/models"
	"carbon-api/utils"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type CarbonSummaryRepository interface {
	GetCarbonSummary(userId int) (models.CarbonSummaryResponse, int, error)
	UpdateCarbonSummary(userId int) (models.CarbonSummary, int, error)
}

type carbonSummaryRepository struct {
	DB *gorm.DB
}

func NewCarbonSummaryRepository(DB *gorm.DB) CarbonSummaryRepository {
	return &carbonSummaryRepository{DB}
}

func (repo *carbonSummaryRepository) GetCarbonSummary(userId int) (models.CarbonSummaryResponse, int, error) {
	var carbonSummary models.CarbonSummaryResponse
	result := repo.DB.Table("carbon_summaries").
		Select("carbon_summaries.user_id, users.name as user_name, users.email as user_email, carbon_summaries.fuel_emission, carbon_summaries.electric_emission, carbon_summaries.total_emission, carbon_summaries.total_tree").
		Joins("JOIN users ON carbon_summaries.user_id = users.id").
		Where("carbon_summaries.user_id = ?", userId).
		Scan(&carbonSummary)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || carbonSummary.UserID == 0 {
		return carbonSummary, http.StatusNotFound, errors.New("Carbon summary not found")
	}

	return carbonSummary, http.StatusOK, nil
}

func (repo *carbonSummaryRepository) UpdateCarbonSummary(userId int) (models.CarbonSummary, int, error) {
	// Calculate fuel emission amount
	carbonFuels, _, err := NewCarbonFuelRepository(repo.DB).GetAllCarbonFuels(userId)
	if err != nil {
		return models.CarbonSummary{}, http.StatusInternalServerError, err
	}

	var fuelEmission float64
	for _, carbonFuel := range carbonFuels {
		fuelEmission += utils.CalculateEmissionAmount(carbonFuel.TotalConsumption, carbonFuel.EmissionFactor)
	}

	// TODO: Calculate eletric emission amount
	var electricEmission float64

	// Calculate total emission amount & tree
	totalEmission := utils.CalculateTotalEmission(fuelEmission, electricEmission)
	totalTree := utils.CalculateTotalTree(totalEmission)

	// Update carbon summary
	var carbonSummary models.CarbonSummary
	result := repo.DB.Table("carbon_summaries").
		Where("user_id = ?", userId).
		First(&carbonSummary)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || carbonSummary.ID == 0 {
		carbonSummary = models.CarbonSummary{
			UserID:           userId,
			FuelEmission:     fuelEmission,
			ElectricEmission: electricEmission,
			TotalEmission:    totalEmission,
			TotalTree:        totalTree,
		}

		result = repo.DB.Table("carbon_summaries").
			Create(&carbonSummary)
	} else {
		carbonSummary.FuelEmission = fuelEmission
		carbonSummary.ElectricEmission = electricEmission
		carbonSummary.TotalEmission = totalEmission
		carbonSummary.TotalTree = totalTree

		result = repo.DB.Table("carbon_summaries").
			Save(&carbonSummary)
	}

	if result.Error != nil {
		return models.CarbonSummary{}, http.StatusInternalServerError, result.Error
	}

	return carbonSummary, http.StatusOK, nil
}
