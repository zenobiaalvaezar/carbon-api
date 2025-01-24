package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCarbonSummary_Success(t *testing.T) {
	e := echo.New()
	mockRepository := new(repositories.MockCarbonSummaryRepository)

	expectedResponse := models.CarbonSummaryResponse{
		UserID:           1,
		UserName:         "John Doe",
		UserEmail:        "john.doe@example.com",
		FuelEmission:     50.5,
		ElectricEmission: 30.0,
		TotalEmission:    80.5,
		TotalTree:        10,
	}

	mockRepository.On("GetCarbonSummary", 1).Return(expectedResponse, http.StatusOK, nil)

	controller := NewCarbonSummaryController(mockRepository)

	req := httptest.NewRequest(http.MethodGet, "/carbon-summaries", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user_id", 1)

	err := controller.GetCarbonSummary(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualResponse models.CarbonSummaryResponse
	json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	assert.Equal(t, expectedResponse, actualResponse)

	mockRepository.AssertExpectations(t)
}
