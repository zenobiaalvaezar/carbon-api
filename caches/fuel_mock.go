package caches

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type FuelCacheMock struct {
	mock.Mock
}

func (m *FuelCacheMock) GetAllFuels() ([]models.Fuel, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelCacheMock) GetFuelByID(id int) (models.Fuel, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelCacheMock) CreateAllFuels(fuels []models.Fuel) (int, error) {
	args := m.Called(fuels)
	return args.Int(0), args.Error(1)
}

func (m *FuelCacheMock) CreateFuel(fuel models.Fuel) (int, error) {
	args := m.Called(fuel)
	return args.Int(0), args.Error(1)
}

func (m *FuelCacheMock) UpdateFuel(fuel models.Fuel) (int, error) {
	args := m.Called(fuel)
	return args.Int(0), args.Error(1)
}

func (m *FuelCacheMock) DeleteFuel(id int) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
