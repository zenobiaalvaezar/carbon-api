package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CarbonElectricController struct {
	CarbonElectricRepository repositories.CarbonElectricRepository
}

func NewCarbonElectricController(carbonElectricRepository repositories.CarbonElectricRepository) *CarbonElectricController {
	return &CarbonElectricController{carbonElectricRepository}
}

func (ctrl *CarbonElectricController) GetAllCarbonElectrics(c echo.Context) error {
	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing

	carbonElectrics, status, err := ctrl.CarbonElectricRepository.GetAllCarbonElectrics(userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonElectrics)
}

func (ctrl *CarbonElectricController) GetCarbonElectricByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon electric ID"})
	}

	carbonElectric, status, err := ctrl.CarbonElectricRepository.GetCarbonElectricByID(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonElectric)
}

func (ctrl *CarbonElectricController) CreateCarbonElectric(c echo.Context) error {
	var carbonElectricRequest models.CarbonElectricRequest
	if err := c.Bind(&carbonElectricRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if carbonElectricRequest.UsageType != "consumption" && carbonElectricRequest.UsageType != "rupiah" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid usage type"})
	}

	if carbonElectricRequest.UsageAmount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Usage amount must be greater than 0"})
	}

	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing
	carbonElectricRequest.UserID = userId
	carbonElectric, status, err := ctrl.CarbonElectricRepository.CreateCarbonElectric(carbonElectricRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonElectric)
}

func (ctrl *CarbonElectricController) DeleteCarbonElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon electric ID"})
	}

	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing
	status, err := ctrl.CarbonElectricRepository.DeleteCarbonElectric(id, userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete carbon electric"})
}
