package repositories

import (
    "carbon-api/models"

    "github.com/stretchr/testify/mock"
)

type MockCarbonElectricRepository struct {
    mock.Mock
}

func (m *MockCarbonElectricRepository) GetAllCarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error) {
    args := m.Called(userId)
    return args.Get(0).([]models.CarbonElectricResponse), args.Int(1), args.Error(2)
}

func (m *MockCarbonElectricRepository) GetCarbonElectricByID(id int) (models.CarbonElectricResponse, int, error) {
    args := m.Called(id)
    return args.Get(0).(models.CarbonElectricResponse), args.Int(1), args.Error(2)
}

func (m *MockCarbonElectricRepository) CreateCarbonElectric(carbonElectric models.CarbonElectricRequest) (models.CarbonElectric, int, error) {
    args := m.Called(carbonElectric)
    return args.Get(0).(models.CarbonElectric), args.Int(1), args.Error(2)
}

func (m *MockCarbonElectricRepository) DeleteCarbonElectric(id int, userId int) (int, error) {
    args := m.Called(id, userId)
    return args.Int(0), args.Error(1)
}

func (m *MockCarbonElectricRepository) GetLast3CarbonElectrics(userId int) ([]models.CarbonElectricResponse, int, error) {
    args := m.Called(userId)
    return args.Get(0).([]models.CarbonElectricResponse), args.Int(1), args.Error(2)
}
