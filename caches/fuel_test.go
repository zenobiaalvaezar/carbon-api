package caches

import (
	"carbon-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFuels(t *testing.T) {
	mockCache := new(FuelCacheMock)
	mockFuels := []models.Fuel{
		{ID: 1, Category: "Category1", Name: "Fuel1", EmissionFactor: 0.5, Price: 100.0, Unit: "liter"},
		{ID: 2, Category: "Category2", Name: "Fuel2", EmissionFactor: 0.7, Price: 120.0, Unit: "liter"},
	}

	mockCache.On("GetAllFuels").Return(mockFuels, 200, nil)

	fuels, status, err := mockCache.GetAllFuels()

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockFuels, fuels)

	mockCache.AssertExpectations(t)
}

func TestGetFuelByID(t *testing.T) {
	mockCache := new(FuelCacheMock)
	mockFuel := models.Fuel{ID: 1, Category: "Category1", Name: "Fuel1", EmissionFactor: 0.5, Price: 100.0, Unit: "liter"}

	mockCache.On("GetFuelByID", 1).Return(mockFuel, 200, nil)

	fuel, status, err := mockCache.GetFuelByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockFuel, fuel)

	mockCache.AssertExpectations(t)
}

func TestCreateAllFuels(t *testing.T) {
	mockCache := new(FuelCacheMock)
	mockFuels := []models.Fuel{
		{ID: 1, Category: "Category1", Name: "Fuel1", EmissionFactor: 0.5, Price: 100.0, Unit: "liter"},
		{ID: 2, Category: "Category2", Name: "Fuel2", EmissionFactor: 0.7, Price: 120.0, Unit: "liter"},
	}

	mockCache.On("CreateAllFuels", mockFuels).Return(201, nil)

	status, err := mockCache.CreateAllFuels(mockFuels)

	assert.NoError(t, err)
	assert.Equal(t, 201, status)

	mockCache.AssertExpectations(t)
}

func TestCreateFuel(t *testing.T) {
	mockCache := new(FuelCacheMock)
	mockFuel := models.Fuel{ID: 1, Category: "Category1", Name: "Fuel1", EmissionFactor: 0.5, Price: 100.0, Unit: "liter"}

	mockCache.On("CreateFuel", mockFuel).Return(201, nil)

	status, err := mockCache.CreateFuel(mockFuel)

	assert.NoError(t, err)
	assert.Equal(t, 201, status)

	mockCache.AssertExpectations(t)
}

func TestUpdateFuel(t *testing.T) {
	mockCache := new(FuelCacheMock)
	mockFuel := models.Fuel{ID: 1, Category: "Category1", Name: "Fuel1", EmissionFactor: 0.5, Price: 100.0, Unit: "liter"}

	mockCache.On("UpdateFuel", mockFuel).Return(200, nil)

	status, err := mockCache.UpdateFuel(mockFuel)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	mockCache.AssertExpectations(t)
}

func TestDeleteFuel(t *testing.T) {
	mockCache := new(FuelCacheMock)

	mockCache.On("DeleteFuel", 1).Return(200, nil)

	status, err := mockCache.DeleteFuel(1)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	mockCache.AssertExpectations(t)
}
