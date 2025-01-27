package repositories

import (
	"carbon-api/models"
	"carbon-api/utils"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type CarbonElectricRepository interface {
	GetAllCarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error)
	GetCarbonElectricByID(id int) (models.CarbonElectricResponse, int, error)
	CreateCarbonElectric(carbonElectric models.CarbonElectricRequest) (models.CarbonElectric, int, error)
	DeleteCarbonElectric(id int, userId int) (int, error)
	GetLast3CarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error)
}

type carbonElectricRepository struct {
	DB *gorm.DB
}

func NewCarbonElectricRepository(DB *gorm.DB) CarbonElectricRepository {
	return &carbonElectricRepository{DB}
}

func (repo *carbonElectricRepository) GetAllCarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error) {
	var carbonElectrics []models.CarbonElectricResponse
	result := repo.DB.Table("carbon_electrics").
		Select("carbon_electrics.id, carbon_electrics.user_id, users.name as user_name, users.email as user_email, electrics.id as electric_id, electrics.province, electrics.price, 'kWh' unit, carbon_electrics.usage_type, carbon_electrics.usage_amount, carbon_electrics.total_consumption, carbon_electrics.emission_factor, carbon_electrics.emission_amount").
		Joins("JOIN users ON carbon_electrics.user_id = users.id").
		Joins("JOIN electrics ON carbon_electrics.electric_id = electrics.id").
		Where("carbon_electrics.user_id = ?", userId).
		Order("carbon_electrics.id").
		Scan(&carbonElectrics)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(carbonElectrics) == 0 {
		carbonElectrics = []models.CarbonElectricResponse{}
	}

	return carbonElectrics, http.StatusOK, nil
}

func (repo *carbonElectricRepository) GetCarbonElectricByID(id int) (models.CarbonElectricResponse, int, error) {
	var carbonElectric models.CarbonElectricResponse
	result := repo.DB.Table("carbon_electrics").
		Select("carbon_electrics.id, carbon_electrics.user_id, users.name as user_name, users.email as user_email, electrics.id as electric_id, electrics.province, electrics.price, 'kWh' unit, carbon_electrics.usage_type, carbon_electrics.usage_amount, carbon_electrics.total_consumption, carbon_electrics.emission_factor, carbon_electrics.emission_amount").
		Joins("JOIN users ON carbon_electrics.user_id = users.id").
		Joins("JOIN electrics ON carbon_electrics.electric_id = electrics.id").
		Where("carbon_electrics.id = ?", id).
		Scan(&carbonElectric)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || carbonElectric.ID == 0 {
		return carbonElectric, http.StatusNotFound, errors.New("Carbon electric not found")
	}

	return carbonElectric, http.StatusOK, nil
}

func (repo *carbonElectricRepository) CreateCarbonElectric(carbonElectric models.CarbonElectricRequest) (models.CarbonElectric, int, error) {
	// Get electric data
	var electric models.Electric
	res := repo.DB.First(&electric, carbonElectric.ElectricID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return models.CarbonElectric{}, http.StatusNotFound, errors.New("Electric not found")
	}

	// Calculate total consumption and emission amount
	totalConsumption := utils.CalculateTotalConsumption(carbonElectric.UsageType, carbonElectric.UsageAmount, electric.Price)
	emissionAmount := utils.CalculateEmissionAmount(totalConsumption, electric.EmissionFactor)

	newCarbonElectric := models.CarbonElectric{
		UserID:           carbonElectric.UserID,
		ElectricID:       carbonElectric.ElectricID,
		UsageType:        carbonElectric.UsageType,
		UsageAmount:      carbonElectric.UsageAmount,
		TotalConsumption: totalConsumption,
		EmissionFactor:   electric.EmissionFactor,
		EmissionAmount:   emissionAmount,
		Price:            electric.Price,
	}

	result := repo.DB.Create(&newCarbonElectric)
	if result.Error != nil {
		return models.CarbonElectric{}, http.StatusInternalServerError, result.Error
	}

	// Update carbon summary
	_, status, err := NewCarbonSummaryRepository(repo.DB).UpdateCarbonSummary(carbonElectric.UserID)
	if err != nil {
		return models.CarbonElectric{}, status, err
	}

	return newCarbonElectric, http.StatusCreated, nil
}

func (repo *carbonElectricRepository) DeleteCarbonElectric(id int, userId int) (int, error) {
	var carbonElectric models.CarbonElectric
	res := repo.DB.First(&carbonElectric, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Carbon electric not found")
	}

	if carbonElectric.UserID != userId {
		return http.StatusForbidden, errors.New("Unauthorized to delete carbon electric")
	}

	result := repo.DB.Delete(&carbonElectric, id)
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

func (repo *carbonElectricRepository) GetLast3CarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error) {
	var carbonElectrics []models.CarbonElectricResponse
	result := repo.DB.Table("carbon_electrics").
		Select("carbon_electrics.id, carbon_electrics.user_id, users.name as user_name, users.email as user_email, electrics.id as electric_id, electrics.province, electrics.price, 'kWh' unit, carbon_electrics.usage_type, carbon_electrics.usage_amount, carbon_electrics.total_consumption, carbon_electrics.emission_factor, carbon_electrics.emission_amount").
		Joins("JOIN users ON carbon_electrics.user_id = users.id").
		Joins("JOIN electrics ON carbon_electrics.electric_id = electrics.id").
		Where("carbon_electrics.user_id = ?", userId).
		Order("carbon_electrics.id").
		Limit(3).
		Scan(&carbonElectrics)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(carbonElectrics) == 0 {
		carbonElectrics = []models.CarbonElectricResponse{}
	}

	return carbonElectrics, http.StatusOK, nil
}
