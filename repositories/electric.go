package repositories

import (
	"carbon-api/models"
	"errors"

	"gorm.io/gorm"
)

type ElectricRepository interface {
	Create(electric *models.Electric) error
	FindByID(id int) (*models.Electric, error)
	FindAll() ([]models.Electric, error)
	Update(id int, electric *models.Electric) error
	Delete(id int) error
	GetAverageEmission(electricID int) (float64, float64, error)
}

type electricRepository struct {
	db *gorm.DB
}

func NewElectricRepository(db *gorm.DB) ElectricRepository {
	return &electricRepository{db: db}
}

func (r *electricRepository) Create(electric *models.Electric) error {
	return r.db.Create(electric).Error
}

func (repo *electricRepository) FindByID(id int) (*models.Electric, error) {
	var electric models.Electric

	err := repo.db.First(&electric, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &electric, err
}

func (repo *electricRepository) FindAll() ([]models.Electric, error) {
	var electrics []models.Electric

	err := repo.db.Find(&electrics).Error
	return electrics, err
}

func (repo *electricRepository) Update(id int, electric *models.Electric) error {
	var existing models.Electric

	err := repo.db.First(&existing, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}
	return repo.db.Model(&existing).Updates(electric).Error
}

func (repo *electricRepository) Delete(id int) error {
	return repo.db.Delete(&models.Electric{}, id).Error
}

func (repo *electricRepository) GetAverageEmission(electricID int) (float64, float64, error) {
	var nationalAverageEmission float64
	var provincialAverageEmission float64
	var err error

	err = repo.db.Model(&models.Electric{}).Select("AVG(emission_factor)").Scan(&nationalAverageEmission).Error
	if err != nil {
		return 0, 0, err
	}

	if electricID > 0 {
		err = repo.db.Model(&models.Electric{}).Where("id = ?", electricID).
			Select("AVG(emission_factor)").Scan(&provincialAverageEmission).Error
		if err != nil {
			return 0, 0, err
		}
	} else {
		provincialAverageEmission = 0
	}

	return nationalAverageEmission, provincialAverageEmission, nil
}
