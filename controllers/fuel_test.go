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

func TestFuelController_GetAllFuels(t *testing.T) {
	e := echo.New()

	fuelCacheMock := new(caches.FuelCacheMock)
	fuelRepoMock := new(repositories.FuelRepositoryMock)

	fuels := []models.Fuel{
		{ID: 1, Category: "Gas", Name: "Fuel A", EmissionFactor: 2.5, Price: 3.2, Unit: "L"},
		{ID: 2, Category: "Diesel", Name: "Fuel B", EmissionFactor: 2.8, Price: 3.5, Unit: "L"},
	}

	fuelCacheMock.On("GetAllFuels").Return(fuels, http.StatusOK, nil)

	fuelRepoMock.On("GetAllFuels").Return([]models.Fuel{}, http.StatusOK, nil)

	ctrl := NewFuelController(fuelRepoMock, fuelCacheMock)

	req := httptest.NewRequest(http.MethodGet, "/fuels", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := ctrl.GetAllFuels(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"id":1,"category":"Gas","name":"Fuel A","emission_factor":2.5,"price":3.2,"unit":"L"},{"id":2,"category":"Diesel","name":"Fuel B","emission_factor":2.8,"price":3.5,"unit":"L"}]`, rec.Body.String())

	fuelCacheMock.AssertExpectations(t)
}
