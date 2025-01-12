package repositories

import (
	"carbon-api/models"
	"carbon-api/utils"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type CarbonFuelRepository interface {
	GetAllCarbonFuels(userId int) ([]models.CarbonFuelResponse, int, error)
	GetCarbonFuelByID(id int) (models.CarbonFuelResponse, int, error)
	CreateCarbonFuel(carbonFuel models.CarbonFuelRequest) (models.CarbonFuel, int, error)
	DeleteCarbonFuel(id int, userId int) (int, error)
}

type carbonFuelRepository struct {
	DB *gorm.DB
}

func NewCarbonFuelRepository(DB *gorm.DB) CarbonFuelRepository {
	return &carbonFuelRepository{DB}
}

func (repo *carbonFuelRepository) GetAllCarbonFuels(userId int) ([]models.CarbonFuelResponse, int, error) {
	var carbonFuels []models.CarbonFuelResponse
	result := repo.DB.Table("carbon_fuels").
		Select("carbon_fuels.id, carbon_fuels.user_id, users.name as user_name, users.email as user_email, carbon_fuels.fuel_id, fuels.name as fuel_name, carbon_fuels.price, fuels.unit, carbon_fuels.usage_amount, carbon_fuels.usage_type, carbon_fuels.total_consumption, carbon_fuels.emission_factor, carbon_fuels.emission_amount").
		Joins("JOIN users ON carbon_fuels.user_id = users.id").
		Joins("JOIN fuels ON carbon_fuels.fuel_id = fuels.id").
		Where("carbon_fuels.user_id = ?", userId).
		Scan(&carbonFuels)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(carbonFuels) == 0 {
		carbonFuels = []models.CarbonFuelResponse{}
	}

	return carbonFuels, http.StatusOK, nil
}

func (repo *carbonFuelRepository) GetCarbonFuelByID(id int) (models.CarbonFuelResponse, int, error) {
	// Get carbon fuel response
	var carbonFuel models.CarbonFuelResponse
	result := repo.DB.Table("carbon_fuels").
		Select("carbon_fuels.id, carbon_fuels.user_id, users.name as user_name, users.email as user_email, carbon_fuels.fuel_id, fuels.name as fuel_name, carbon_fuels.price, fuels.unit, carbon_fuels.usage_amount, carbon_fuels.usage_type, carbon_fuels.total_consumption, carbon_fuels.emission_factor, carbon_fuels.emission_amount").
		Joins("JOIN users ON carbon_fuels.user_id = users.id").
		Joins("JOIN fuels ON carbon_fuels.fuel_id = fuels.id").
		Where("carbon_fuels.id = ?", id).
		Scan(&carbonFuel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || carbonFuel.ID == 0 {
		return carbonFuel, http.StatusNotFound, errors.New("Carbon fuel not found")
	}

	return carbonFuel, http.StatusOK, nil
}

func (repo *carbonFuelRepository) CreateCarbonFuel(carbonFuel models.CarbonFuelRequest) (models.CarbonFuel, int, error) {
	// Get fuel data
	var fuel models.Fuel
	res := repo.DB.First(&fuel, carbonFuel.FuelID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return models.CarbonFuel{}, http.StatusNotFound, errors.New("Fuel not found")
	}

	// Create new carbon fuel
	totalConsumption := utils.CalculateTotalConsumption(carbonFuel.UsageType, carbonFuel.UsageAmount, fuel.Price)
	emissionAmount := utils.CalculateEmissionAmount(totalConsumption, fuel.EmissionFactor)

	newCarbonFuel := models.CarbonFuel{
		UserID:           carbonFuel.UserID,
		FuelID:           carbonFuel.FuelID,
		Price:            fuel.Price,
		UsageAmount:      carbonFuel.UsageAmount,
		UsageType:        carbonFuel.UsageType,
		TotalConsumption: totalConsumption,
		EmissionFactor:   fuel.EmissionFactor,
		EmissionAmount:   emissionAmount,
	}

	result := repo.DB.Create(&newCarbonFuel)
	if result.Error != nil {
		return models.CarbonFuel{}, http.StatusInternalServerError, result.Error
	}

	// Update carbon summary
	_, status, err := NewCarbonSummaryRepository(repo.DB).UpdateCarbonSummary(carbonFuel.UserID)
	if err != nil {
		return models.CarbonFuel{}, status, err
	}

	return newCarbonFuel, http.StatusCreated, nil
}

func (repo *carbonFuelRepository) DeleteCarbonFuel(id int, userId int) (int, error) {
	// Check if carbon fuel exists
	var carbonFuel models.CarbonFuel
	res := repo.DB.First(&carbonFuel, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Carbon fuel not found")
	}

	// Check if user is authorized to delete the carbon fuel
	if carbonFuel.UserID != userId {
		return http.StatusForbidden, errors.New("Unauthorized to delete carbon fuel")
	}

	// Delete carbon fuel
	result := repo.DB.Delete(&carbonFuel, id)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	// Update carbon summary
	_, status, err := NewCarbonSummaryRepository(repo.DB).UpdateCarbonSummary(userId)
	if err != nil {
		return status, err
	}

	return http.StatusOK, nil
}
