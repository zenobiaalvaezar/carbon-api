package controllers

import (
	"carbon-api/caches"
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllElectrics(t *testing.T) {
	e := echo.New()

	mockCache := new(caches.MockElectricCache)
	mockRepository := new(repositories.MockElectricRepository)

	electrics := []models.Electric{
		{ID: 1, Province: "Province A", EmissionFactor: 0.5, Price: 100.0},
		{ID: 2, Province: "Province B", EmissionFactor: 0.7, Price: 150.0},
	}
	mockCache.On("GetAllElectrics").Return(electrics, http.StatusOK, nil)
	mockRepository.On("FindAll").Return(electrics, nil)

	controller := NewElectricController(mockRepository, mockCache)

	req := httptest.NewRequest(http.MethodGet, "/electrics", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetAllElectrics(ctx)
	if err != nil {
		t.Fatalf("Error calling GetAllElectrics: %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody := `[{"id":1,"province":"Province A","emission_factor":0.5,"price":100.0},{"id":2,"province":"Province B","emission_factor":0.7,"price":150.0}]`
	assert.JSONEq(t, expectedBody, rec.Body.String())
	mockCache.AssertExpectations(t) //
}

func TestGetElectricByID(t *testing.T) {
	mockRepo := new(repositories.MockElectricRepository)
	mockCache := new(caches.MockElectricCache)

	electric := &models.Electric{ID: 1, Province: "Province1", EmissionFactor: 0.5, Price: 100}

	mockRepo.On("FindByID", 1).Return(electric, nil)
	mockCache.On("CreateElectric", *electric).Return(http.StatusOK, nil)

	ctrl := NewElectricController(mockRepo, mockCache)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/electrics/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, ctrl.GetElectricByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"province":"Province1"`)
	}

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestDeleteElectric(t *testing.T) {
	mockRepo := new(repositories.MockElectricRepository)
	mockCache := new(caches.MockElectricCache)

	mockRepo.On("Delete", 1).Return(nil)
	mockCache.On("DeleteElectric", 1).Return(http.StatusOK, nil)

	ctrl := NewElectricController(mockRepo, mockCache)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/electrics/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, ctrl.DeleteElectric(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Electric deleted successfully")
	}

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
