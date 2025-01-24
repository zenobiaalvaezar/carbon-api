package caches

import (
	"carbon-api/models"

	"github.com/stretchr/testify/mock"
)

type TreeCacheMock struct {
	mock.Mock
}

func (m *TreeCacheMock) GetAllTrees() ([]models.Tree, int, error) {
	args := m.Called()
	return args.Get(0).([]models.Tree), args.Int(1), args.Error(2)
}

func (m *TreeCacheMock) GetTreeByID(id int) (models.Tree, int, error) {
	args := m.Called(id)
	return args.Get(0).(models.Tree), args.Int(1), args.Error(2)
}

func (m *TreeCacheMock) SetAllTrees(trees []models.Tree) error {
	args := m.Called(trees)
	return args.Error(0)
}

func (m *TreeCacheMock) SetTree(tree models.Tree) error {
	args := m.Called(tree)
	return args.Error(0)
}

func (m *TreeCacheMock) DeleteTree(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
