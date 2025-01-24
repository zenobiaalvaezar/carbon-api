package repositories

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type MockTreeRepository struct {
	mock.Mock
}

func (m *MockTreeRepository) GetAllTrees() ([]models.Tree, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Tree), args.Int(1), args.Error(2)
}

func (m *MockTreeRepository) GetTreeByID(id int) (models.Tree, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Tree), args.Int(1), args.Error(2)
}

func (m *MockTreeRepository) CreateTree(tree *models.Tree) (int, error) {
	args := m.Called(tree)
	return args.Int(0), args.Error(1)
}

func (m *MockTreeRepository) UpdateTree(tree *models.Tree) (int, error) {
	args := m.Called(tree)
	return args.Int(0), args.Error(1)
}

func (m *MockTreeRepository) DeleteTree(id int) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}

type MockTreeCache struct {
	mock.Mock
}

func (m *MockTreeCache) GetAllTrees() ([]models.Tree, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Tree), args.Int(1), args.Error(2)
}

func (m *MockTreeCache) GetTreeByID(id int) (models.Tree, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Tree), args.Int(1), args.Error(2)
}

func (m *MockTreeCache) SetAllTrees(trees []models.Tree) error {
	args := m.Called(trees)
	return args.Error(0)
}

func (m *MockTreeCache) SetTree(tree models.Tree) error {
	args := m.Called(tree)
	return args.Error(0)
}

func (m *MockTreeCache) DeleteTree(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
