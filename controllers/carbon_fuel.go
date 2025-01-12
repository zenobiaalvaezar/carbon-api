package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CarbonFuelController struct {
	CarbonFuelRepository repositories.CarbonFuelRepository
}

func NewCarbonFuelController(carbonFuelRepository repositories.CarbonFuelRepository) *CarbonFuelController {
	return &CarbonFuelController{carbonFuelRepository}
}

func (ctrl *CarbonFuelController) GetAllCarbonFuels(c echo.Context) error {
	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing

	carbonFuels, status, err := ctrl.CarbonFuelRepository.GetAllCarbonFuels(userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonFuels)
}

func (ctrl *CarbonFuelController) GetCarbonFuelByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon fuel ID"})
	}

	carbonFuel, status, err := ctrl.CarbonFuelRepository.GetCarbonFuelByID(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonFuel)
}

func (ctrl *CarbonFuelController) CreateCarbonFuel(c echo.Context) error {
	var carbonFuelRequest models.CarbonFuelRequest
	if err := c.Bind(&carbonFuelRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if carbonFuelRequest.UsageType != "consumption" && carbonFuelRequest.UsageType != "rupiah" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid usage type"})
	}

	if carbonFuelRequest.UsageAmount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Usage amount must be greater than 0"})
	}

	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing
	carbonFuelRequest.UserID = userId
	carbonFuel, status, err := ctrl.CarbonFuelRepository.CreateCarbonFuel(carbonFuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonFuel)
}

func (ctrl *CarbonFuelController) DeleteCarbonFuel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon fuel ID"})
	}

	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing
	status, err := ctrl.CarbonFuelRepository.DeleteCarbonFuel(id, userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete carbon fuel"})
}
