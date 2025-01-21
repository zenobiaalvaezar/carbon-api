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
