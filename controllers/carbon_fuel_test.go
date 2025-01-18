package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCarbonFuels(t *testing.T) {
	mockRepo := new(repositories.MockCarbonFuelRepository)
	ctrl := NewCarbonFuelController(mockRepo)

	expectedFuels := []models.CarbonFuelResponse{
		{
			ID:               1,
			UserID:           2,
			FuelID:           1,
			FuelName:         "Fuel A",
			Price:            100,
			Unit:             "kg",
			UsageAmount:      10,
			UsageType:        "consumption",
			TotalConsumption: 100,
			EmissionFactor:   1.5,
			EmissionAmount:   150,
			UserEmail:        "",
			UserName:         "",
		},
	}
	mockRepo.On("GetAllCarbonFuels", 2).Return(expectedFuels, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/carbon-fuels", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := ctrl.GetAllCarbonFuels(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"id":1,"user_id":2,"fuel_id":1,"fuel_name":"Fuel A","price":100,"unit":"kg","usage_amount":10,"usage_type":"consumption","total_consumption":100,"emission_factor":1.5,"emission_amount":150,"user_email":"","user_name":""}]`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}
