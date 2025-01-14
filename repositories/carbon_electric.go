package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type CarbonElectricRepository interface {
	GetAllCarbonElectrics(userId int) ([]models.CarbonElectric, int, error)
	CreateCarbonElectric(carbonElectric models.CarbonElectric) (models.CarbonElectric, int, error)
	SoftDeleteCarbonElectric(id int, userId int) (int, error)
}

type carbonElectricRepository struct {
	DB *gorm.DB
}

func NewCarbonElectricRepository(DB *gorm.DB) CarbonElectricRepository {
	return &carbonElectricRepository{DB}
}

func (repo *carbonElectricRepository) GetAllCarbonElectrics(userId int) ([]models.CarbonElectric, int, error) {
	var carbonElectrics []models.CarbonElectric
	result := repo.DB.Where("user_id = ?", userId).Find(&carbonElectrics)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(carbonElectrics) == 0 {
		carbonElectrics = []models.CarbonElectric{}
	}

	return carbonElectrics, http.StatusOK, nil
}

func (repo *carbonElectricRepository) CreateCarbonElectric(carbonElectric models.CarbonElectric) (models.CarbonElectric, int, error) {
	result := repo.DB.Create(&carbonElectric)
	if result.Error != nil {
		return models.CarbonElectric{}, http.StatusInternalServerError, result.Error
	}
	return carbonElectric, http.StatusCreated, nil
}

func (repo *carbonElectricRepository) SoftDeleteCarbonElectric(id int, userId int) (int, error) {
	var carbonElectric models.CarbonElectric
	result := repo.DB.First(&carbonElectric, "id = ? AND user_id = ?", id, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Carbon electric record not found")
	}

	if err := repo.DB.Delete(&carbonElectric).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
