package controllers

import (
	"carbon-api/caches"
	"carbon-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTrees(t *testing.T) {
	mockCache := new(caches.TreeCacheMock)
	mockTrees := []models.Tree{
		{ID: 1, TreeCategoryID: 1, Name: "Oak", Description: "Oak Tree", Price: 100.0, Stock: 50},
		{ID: 2, TreeCategoryID: 2, Name: "Pine", Description: "Pine Tree", Price: 150.0, Stock: 30},
	}

	// Set up the mock to return mockTrees with HTTP status OK
	mockCache.On("GetAllTrees").Return(mockTrees, 200, nil)

	// Call the method you're testing
	trees, status, err := mockCache.GetAllTrees()

	// Assert the expected results
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockTrees, trees)

	// Assert that the mock method was called
	mockCache.AssertExpectations(t)
}