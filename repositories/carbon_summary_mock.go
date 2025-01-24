package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockCarbonSummaryRepository struct {
	mock.Mock
}

func (m *MockCarbonSummaryRepository) GetCarbonSummary(userId int) (models.CarbonSummaryResponse, int, error) {
	args := m.Called(userId)
	return args.Get(0).(models.CarbonSummaryResponse), args.Int(1), args.Error(2)
}

func (m *MockCarbonSummaryRepository) UpdateCarbonSummary(userId int) (models.CarbonSummary, int, error) {
	args := m.Called(userId)
	return args.Get(0).(models.CarbonSummary), args.Int(1), args.Error(2)
}
