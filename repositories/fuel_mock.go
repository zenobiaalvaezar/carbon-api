package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type FuelRepositoryMock struct {
	mock.Mock
}

func (m *FuelRepositoryMock) GetAllFuels() ([]models.Fuel, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelRepositoryMock) GetFuelByID(id int) (models.Fuel, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelRepositoryMock) CreateFuel(fuel models.FuelRequest) (models.Fuel, int, error) {
	args := m.Called(fuel)
	return args.Get(0).(models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelRepositoryMock) UpdateFuel(id int, fuel models.FuelRequest) (models.Fuel, int, error) {
	args := m.Called(id, fuel)
	return args.Get(0).(models.Fuel), args.Int(1), args.Error(2)
}

func (m *FuelRepositoryMock) DeleteFuel(id int) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
