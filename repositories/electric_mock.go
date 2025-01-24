package repositories

import (
        "carbon-api/models"

        "github.com/stretchr/testify/mock"
)

type MockElectricRepository struct {
        mock.Mock
}

func (m *MockElectricRepository) Create(electric *models.Electric) error {
        args := m.Called(electric)
        return args.Error(0)
}

func (m *MockElectricRepository) FindByID(id int) (*models.Electric, error) {
        args := m.Called(id)
        return args.Get(0).(*models.Electric), args.Error(1)
}

func (m *MockElectricRepository) FindAll() ([]models.Electric, error) {
        args := m.Called()
        return args.Get(0).([]models.Electric), args.Error(1)
}

func (m *MockElectricRepository) Update(id int, electric *models.Electric) error {
        args := m.Called(id, electric)
        return args.Error(0)
}

func (m *MockElectricRepository) Delete(id int) error {
        args := m.Called(id)
        return args.Error(0)
}

func (m *MockElectricRepository) GetAverageEmission(electricID int) (float64, float64, error) {
        args := m.Called(electricID)
        return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}