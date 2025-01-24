package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockCarbonFuelRepository struct {
	mock.Mock
}

func (m *MockCarbonFuelRepository) GetAllCarbonFuels(userId int) ([]models.CarbonFuelResponse, int, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.CarbonFuelResponse), args.Int(1), args.Error(2)
}

func (m *MockCarbonFuelRepository) GetCarbonFuelByID(id int) (models.CarbonFuelResponse, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.CarbonFuelResponse), args.Int(1), args.Error(2)
}

func (m *MockCarbonFuelRepository) CreateCarbonFuel(carbonFuel models.CarbonFuelRequest) (models.CarbonFuel, int, error) {
	args := m.Called(carbonFuel)
	return args.Get(0).(models.CarbonFuel), args.Int(1), args.Error(2)
}

func (m *MockCarbonFuelRepository) DeleteCarbonFuel(id int, userId int) (int, error) {
	args := m.Called(id, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockCarbonFuelRepository) GetLast3CarbonFuels(userId int) ([]models.CarbonFuelResponse, int, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.CarbonFuelResponse), args.Int(1), args.Error(2)
}
