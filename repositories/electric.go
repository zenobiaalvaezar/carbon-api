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
