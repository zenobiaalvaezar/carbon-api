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

func TestGetAllCarbonElectrics_Success(t *testing.T) {
	e := echo.New()

	mockRepository := new(repositories.MockCarbonElectricRepository)

	carbonElectrics := []models.CarbonElectricResponse{
		{ID: 1, UserID: 2, UserName: "John Doe", UserEmail: "john@example.com", ElectricID: 1, Province: "DKI Jakarta", Price: 1500, Unit: "kWh", UsageType: "consumption", UsageAmount: 100.0, TotalConsumption: 100.0, EmissionFactor: 0.5, EmissionAmount: 50.0},
		{ID: 2, UserID: 2, UserName: "Jane Smith", UserEmail: "jane@example.com", ElectricID: 1, Province: "DKI Jakarta", Price: 1500, Unit: "kWh", UsageType: "rupiah", UsageAmount: 150.0, TotalConsumption: 150.0, EmissionFactor: 0.6, EmissionAmount: 90.0},
	}
	mockRepository.On("GetAllCarbonElectrics", 2).Return(carbonElectrics, http.StatusOK, nil)

	controller := NewCarbonElectricController(mockRepository)

	req := httptest.NewRequest(http.MethodGet, "/carbon-electrics", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user_id", 2)

	err := controller.GetAllCarbonElectrics(ctx)
	if err != nil {
		t.Fatalf("Error calling GetAllCarbonElectrics: %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody := `[{"id":1,"user_id":2,"user_name":"John Doe","user_email":"john@example.com","electric_id":1,"province":"DKI Jakarta","price":1500,"unit":"kWh","usage_type":"consumption","usage_amount":100.0,"total_consumption":100.0,"emission_factor":0.5,"emission_amount":50.0},{"id":2,"user_id":2,"user_name":"Jane Smith","user_email":"jane@example.com","electric_id":1,"province":"DKI Jakarta","price":1500,"unit":"kWh","usage_type":"rupiah","usage_amount":150.0,"total_consumption":150.0,"emission_factor":0.6,"emission_amount":90.0}]`
	assert.JSONEq(t, expectedBody, rec.Body.String())

	mockRepository.AssertExpectations(t)
}

func TestGetAllCarbonElectrics_Failure(t *testing.T) {
	e := echo.New()

	mockRepository := new(repositories.MockCarbonElectricRepository)

	mockRepository.On("GetAllCarbonElectrics", 2).Return([]models.CarbonElectricResponse{}, http.StatusOK, nil)

	controller := NewCarbonElectricController(mockRepository)

	req := httptest.NewRequest(http.MethodGet, "/carbon-electrics", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user_id", 2)

	err := controller.GetAllCarbonElectrics(ctx)
	if err != nil {
		t.Fatalf("Error calling GetAllCarbonElectrics (empty list): %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedEmptyBody := `[]`
	assert.JSONEq(t, expectedEmptyBody, rec.Body.String())

	mockRepository.AssertExpectations(t)
}

func TestGetCarbonElectricByID_Success(t *testing.T) {
	e := echo.New()

	mockRepository := new(repositories.MockCarbonElectricRepository)

	carbonElectric := models.CarbonElectricResponse{
		ID:               1,
		UserID:           2,
		UserName:         "User A",
		UserEmail:        "usera@example.com",
		ElectricID:       1,
		Province:         "DKI Jakarta",
		Price:            1500,
		Unit:             "kWh",
		UsageType:        "consumption",
		UsageAmount:      100,
		TotalConsumption: 200,
		EmissionFactor:   0.5,
		EmissionAmount:   100,
	}
	mockRepository.On("GetCarbonElectricByID", 1).Return(carbonElectric, http.StatusOK, nil)

	controller := NewCarbonElectricController(mockRepository)

	ce := e.Group("/carbon-electrics")
	ce.GET("/:id", controller.GetCarbonElectricByID)

	req := httptest.NewRequest(http.MethodGet, "/carbon-electrics/1", nil) // Valid ID "1"
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	ctx.Set("user_id", 2)

	err := controller.GetCarbonElectricByID(ctx)
	if err != nil {
		t.Fatalf("Error calling GetCarbonElectricByID: %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `{"id":1,"user_id":2,"user_name":"User A","user_email":"usera@example.com","electric_id":1,"province":"DKI Jakarta","price":1500,"unit":"kWh","usage_type":"consumption","usage_amount":100,"total_consumption":200,"emission_factor":0.5,"emission_amount":100}`
	assert.JSONEq(t, expectedBody, rec.Body.String())

	mockRepository.AssertExpectations(t)
}

func TestGetCarbonElectricByID_Failure_InvalidID(t *testing.T) {
	e := echo.New()

	mockRepository := new(repositories.MockCarbonElectricRepository)

	controller := NewCarbonElectricController(mockRepository)

	req := httptest.NewRequest(http.MethodGet, "/carbon-electrics/invalid", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user_id", 2)

	err := controller.GetCarbonElectricByID(ctx)
	if err != nil {
		t.Fatalf("Error calling GetCarbonElectricByID (invalid ID): %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expectedBody := `{"message":"Invalid carbon electric ID"}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}
