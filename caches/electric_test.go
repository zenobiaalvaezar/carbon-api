package caches

import (
	"carbon-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllElectrics(t *testing.T) {
	mockCache := new(MockElectricCache)
	mockElectrics := []models.Electric{
		{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100.0},
		{ID: 2, Province: "Province2", EmissionFactor: 0.7, Price: 120.0},
	}

	mockCache.On("GetAllElectrics").Return(mockElectrics, 200, nil)

	electrics, status, err := mockCache.GetAllElectrics()

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockElectrics, electrics)

	mockCache.AssertExpectations(t)
}

func TestGetElectricByID(t *testing.T) {
	mockCache := new(MockElectricCache)
	mockElectric := models.Electric{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100.0}

	mockCache.On("GetElectricByID", 1).Return(mockElectric, 200, nil)

	electric, status, err := mockCache.GetElectricByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, mockElectric, electric)

	mockCache.AssertExpectations(t)
}

func TestCreateAllElectrics(t *testing.T) {
	mockCache := new(MockElectricCache)
	mockElectrics := []models.Electric{
		{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100.0},
		{ID: 2, Province: "Province2", EmissionFactor: 0.7, Price: 120.0},
	}

	mockCache.On("CreateAllElectrics", mockElectrics).Return(201, nil)

	status, err := mockCache.CreateAllElectrics(mockElectrics)

	assert.NoError(t, err)
	assert.Equal(t, 201, status)

	mockCache.AssertExpectations(t)
}

func TestCreateElectric(t *testing.T) {
	mockCache := new(MockElectricCache)
	mockElectric := models.Electric{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100.0}

	mockCache.On("CreateElectric", mockElectric).Return(201, nil)

	status, err := mockCache.CreateElectric(mockElectric)

	assert.NoError(t, err)
	assert.Equal(t, 201, status)

	mockCache.AssertExpectations(t)
}

func TestUpdateElectric(t *testing.T) {
	mockCache := new(MockElectricCache)
	mockElectric := models.Electric{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100.0}

	mockCache.On("UpdateElectric", mockElectric).Return(200, nil)

	status, err := mockCache.UpdateElectric(mockElectric)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	mockCache.AssertExpectations(t)
}

func TestDeleteElectric(t *testing.T) {
	mockCache := new(MockElectricCache)

	mockCache.On("DeleteElectric", 1).Return(200, nil)

	status, err := mockCache.DeleteElectric(1)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	mockCache.AssertExpectations(t)
}
