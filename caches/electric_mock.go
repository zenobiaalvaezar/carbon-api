package caches

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockElectricCache struct {
	mock.Mock
}

func (m *MockElectricCache) GetAllElectrics() ([]models.Electric, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Electric), args.Int(1), args.Error(2)
}

func (m *MockElectricCache) GetElectricByID(id int) (models.Electric, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Electric), args.Int(1), args.Error(2)
}

func (m *MockElectricCache) CreateAllElectrics(electrics []models.Electric) (int, error) {
	args := m.Called(electrics)
	return args.Int(0), args.Error(1)
}

func (m *MockElectricCache) CreateElectric(electric models.Electric) (int, error) {
	args := m.Called(electric)
	return args.Int(0), args.Error(1)
}

func (m *MockElectricCache) UpdateElectric(electric models.Electric) (int, error) {
	args := m.Called(electric)
	return args.Int(0), args.Error(1)
}

func (m *MockElectricCache) DeleteElectric(id int) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
