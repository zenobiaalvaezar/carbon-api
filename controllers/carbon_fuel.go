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

// GetAllCarbonFuels godoc
// @Summary Get all carbon fuels
// @Description Get all carbon fuels for a specific user
// @Tags CarbonFuels
// @Accept json
// @Produce json
// @Success 200 {array} models.CarbonFuelResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carbon-fuels [get]
func (ctrl *CarbonFuelController) GetAllCarbonFuels(c echo.Context) error {
	userId := c.Get("user_id").(int)

	carbonFuels, status, err := ctrl.CarbonFuelRepository.GetAllCarbonFuels(userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonFuels)
}

// GetCarbonFuelByID godoc
// @Summary Get carbon fuel by ID
// @Description Get a specific carbon fuel by its ID
// @Tags CarbonFuels
// @Accept json
// @Produce json
// @Param id path int true "Carbon Fuel ID"
// @Success 200 {object} models.CarbonFuelResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /carbon-fuels/{id} [get]
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

// CreateCarbonFuel godoc
// @Summary Create a new carbon fuel entry
// @Description Create a new carbon fuel entry for a user with specific fuel data
// @Tags CarbonFuels
// @Accept json
// @Produce json
// @Param body body models.CarbonFuelRequest true "Create Carbon Fuel Request"
// @Success 201 {object} models.CarbonFuelResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carbon-fuels [post]
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

	userId := c.Get("user_id").(int)
	carbonFuelRequest.UserID = userId
	carbonFuel, status, err := ctrl.CarbonFuelRepository.CreateCarbonFuel(carbonFuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonFuel)
}

// DeleteCarbonFuel godoc
// @Summary Delete a carbon fuel entry
// @Description Delete a carbon fuel entry by ID for a specific user
// @Tags CarbonFuels
// @Accept json
// @Produce json
// @Param id path int true "Carbon Fuel ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /carbon-fuels/{id} [delete]
func (ctrl *CarbonFuelController) DeleteCarbonFuel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon fuel ID"})
	}

	userId := c.Get("user_id").(int)
	status, err := ctrl.CarbonFuelRepository.DeleteCarbonFuel(id, userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete carbon fuel"})
}
