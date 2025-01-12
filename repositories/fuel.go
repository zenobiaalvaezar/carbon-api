package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type FuelRepository interface {
	GetAllFuels() ([]models.Fuel, int, error)
	GetFuelByID(id int) (models.Fuel, int, error)
	CreateFuel(fuel models.FuelRequest) (models.Fuel, int, error)
	UpdateFuel(id int, fuel models.FuelRequest) (models.Fuel, int, error)
	DeleteFuel(id int) (int, error)
}

type fuelRepository struct {
	DB *gorm.DB
}

func NewFuelRepository(DB *gorm.DB) FuelRepository {
	return &fuelRepository{DB}
}

func (repo *fuelRepository) GetAllFuels() ([]models.Fuel, int, error) {
	var fuels []models.Fuel
	result := repo.DB.Order("id asc").Find(&fuels)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	if len(fuels) == 0 {
		fuels = []models.Fuel{}
	}

	return fuels, http.StatusOK, nil
}

func (repo *fuelRepository) GetFuelByID(id int) (models.Fuel, int, error) {
	var fuel models.Fuel
	result := repo.DB.First(&fuel, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || fuel.ID == 0 {
		return fuel, http.StatusNotFound, errors.New("Fuel not found")
	}

	return fuel, http.StatusOK, nil
}

func (repo *fuelRepository) CreateFuel(fuel models.FuelRequest) (models.Fuel, int, error) {
	newFuel := models.Fuel{
		Category:       fuel.Category,
		Name:           fuel.Name,
		EmissionFactor: fuel.EmissionFactor,
		Price:          fuel.Price,
		Unit:           fuel.Unit,
	}

	result := repo.DB.Create(&newFuel)
	if result.Error != nil {
		return models.Fuel{}, http.StatusInternalServerError, result.Error
	}

	return newFuel, http.StatusCreated, nil
}

func (repo *fuelRepository) UpdateFuel(id int, fuel models.FuelRequest) (models.Fuel, int, error) {
	var existingFuel models.Fuel
	result := repo.DB.First(&existingFuel, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || existingFuel.ID == 0 {
		return existingFuel, http.StatusNotFound, errors.New("Fuel not found")
	}

	existingFuel.Category = fuel.Category
	existingFuel.Name = fuel.Name
	existingFuel.EmissionFactor = fuel.EmissionFactor
	existingFuel.Price = fuel.Price
	existingFuel.Unit = fuel.Unit

	result = repo.DB.Save(&existingFuel)
	if result.Error != nil {
		return models.Fuel{}, http.StatusInternalServerError, result.Error
	}

	return existingFuel, http.StatusOK, nil
}

func (repo *fuelRepository) DeleteFuel(id int) (int, error) {
	result := repo.DB.Delete(&models.Fuel{}, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Fuel not found")
	}

	return http.StatusOK, nil
}
