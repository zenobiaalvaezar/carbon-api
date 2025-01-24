package caches

import (
	"carbon-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTrees(t *testing.T) {
	mockCache := new(TreeCacheMock)
	mockTrees := []models.Tree{
		{ID: 1, TreeCategoryID: 1, Name: "Oak", Description: "Strong tree", Price: 150.50, Stock: 10},
		{ID: 2, TreeCategoryID: 2, Name: "Pine", Description: "Tall tree", Price: 100.75, Stock: 20},
	}

	mockCache.On("GetAllTrees").Return(mockTrees, 200, nil)

	trees, status, err := mockCache.GetAllTrees()

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockTrees, trees)

	mockCache.AssertExpectations(t)
}

func TestGetTreeByID(t *testing.T) {
	mockCache := new(TreeCacheMock)
	mockTree := models.Tree{
		ID: 1, TreeCategoryID: 1, Name: "Oak", Description: "Strong tree", Price: 150.50, Stock: 10,
	}

	mockCache.On("GetTreeByID", 1).Return(mockTree, 200, nil)

	tree, status, err := mockCache.GetTreeByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockTree, tree)

	mockCache.AssertExpectations(t)
}

func TestSetAllTrees(t *testing.T) {
	mockCache := new(TreeCacheMock)
	mockTrees := []models.Tree{
		{ID: 1, TreeCategoryID: 1, Name: "Oak", Description: "Strong tree", Price: 150.50, Stock: 10},
		{ID: 2, TreeCategoryID: 2, Name: "Pine", Description: "Tall tree", Price: 100.75, Stock: 20},
	}

	mockCache.On("SetAllTrees", mockTrees).Return(nil)

	err := mockCache.SetAllTrees(mockTrees)

	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestSetTree(t *testing.T) {
	mockCache := new(TreeCacheMock)
	mockTree := models.Tree{
		ID: 1, TreeCategoryID: 1, Name: "Oak", Description: "Strong tree", Price: 150.50, Stock: 10,
	}

	mockCache.On("SetTree", mockTree).Return(nil)

	err := mockCache.SetTree(mockTree)

	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestDeleteTree(t *testing.T) {
	mockCache := new(TreeCacheMock)

	mockCache.On("DeleteTree", 1).Return(nil)

	err := mockCache.DeleteTree(1)

	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}
