package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"net/http/httptest"
	"strings"
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
	c.Set("user_id", 2)

	err := ctrl.GetAllCarbonFuels(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"id":1,"user_id":2,"fuel_id":1,"fuel_name":"Fuel A","price":100,"unit":"kg","usage_amount":10,"usage_type":"consumption","total_consumption":100,"emission_factor":1.5,"emission_amount":150,"user_email":"","user_name":""}]`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}

func TestGetCarbonFuelByID(t *testing.T) {
	mockRepo := new(repositories.MockCarbonFuelRepository)
	ctrl := NewCarbonFuelController(mockRepo)

	expectedFuel := models.CarbonFuelResponse{
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
	}
	mockRepo.On("GetCarbonFuelByID", 1).Return(expectedFuel, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/carbon-fuels/1", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := ctrl.GetCarbonFuelByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"id":1,"user_id":2,"fuel_id":1,"fuel_name":"Fuel A","price":100,"unit":"kg","usage_amount":10,"usage_type":"consumption","total_consumption":100,"emission_factor":1.5,"emission_amount":150,"user_email":"","user_name":""}`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}

func TestCreateCarbonFuel(t *testing.T) {
	mockRepo := new(repositories.MockCarbonFuelRepository)
	ctrl := NewCarbonFuelController(mockRepo)

	input := models.CarbonFuelRequest{
		UserID:      2,
		FuelID:      1,
		UsageType:   "consumption",
		UsageAmount: 10,
	}
	expectedFuel := models.CarbonFuel{
		ID:               1,
		UserID:           2,
		FuelID:           1,
		Price:            100,
		UsageAmount:      10,
		UsageType:        "consumption",
		TotalConsumption: 100,
		EmissionFactor:   1.5,
		EmissionAmount:   150,
	}
	mockRepo.On("CreateCarbonFuel", input).Return(expectedFuel, http.StatusCreated, nil)

	// Use the JSON body directly in the NewRequest function
	req := httptest.NewRequest(http.MethodPost, "/carbon-fuels", strings.NewReader(`{"fuel_id":1,"usage_type":"consumption","usage_amount":10}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.Set("user_id", 2)

	err := ctrl.CreateCarbonFuel(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"id":1,"user_id":2,"fuel_id":1,"price":100,"usage_amount":10,"usage_type":"consumption","total_consumption":100,"emission_factor":1.5,"emission_amount":150}`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}

func TestDeleteCarbonFuel(t *testing.T) {
	mockRepo := new(repositories.MockCarbonFuelRepository)
	ctrl := NewCarbonFuelController(mockRepo)

	mockRepo.On("DeleteCarbonFuel", 1, 2).Return(http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodDelete, "/carbon-fuels/1", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user_id", 2)

	err := ctrl.DeleteCarbonFuel(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Success delete carbon fuel"}`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}
