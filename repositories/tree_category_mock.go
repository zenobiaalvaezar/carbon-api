package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockTreeCategoryRepository struct {
	mock.Mock
}

func (m *MockTreeCategoryRepository) GetAllTreeCategories() ([]models.TreeCategory, error) {
	args := m.Called()
	return args.Get(0).([]models.TreeCategory), args.Error(1)
}

func (m *MockTreeCategoryRepository) GetTreeCategoryByID(id int) (models.TreeCategory, error) {
	args := m.Called(id)
	return args.Get(0).(models.TreeCategory), args.Error(1)
}

func (m *MockTreeCategoryRepository) CreateTreeCategory(category *models.TreeCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockTreeCategoryRepository) UpdateTreeCategory(category *models.TreeCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockTreeCategoryRepository) DeleteTreeCategory(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
