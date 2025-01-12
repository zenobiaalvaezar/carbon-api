package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FuelController struct {
	FuelRepository repositories.FuelRepository
}

func NewFuelController(fuelRepository repositories.FuelRepository) *FuelController {
	return &FuelController{fuelRepository}
}

func (ctrl *FuelController) GetAllFuels(c echo.Context) error {
	fuels, status, err := ctrl.FuelRepository.GetAllFuels()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, fuels)
}

func (ctrl *FuelController) GetFuelByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid fuel ID"})
	}

	fuel, status, err := ctrl.FuelRepository.GetFuelByID(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, fuel)
}

func (ctrl *FuelController) CreateFuel(c echo.Context) error {
	var fuelRequest models.FuelRequest
	if err := c.Bind(&fuelRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if fuelRequest.Category == "" || fuelRequest.Name == "" || fuelRequest.EmissionFactor <= 0 || fuelRequest.Price <= 0 || fuelRequest.Unit == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	fuel, status, err := ctrl.FuelRepository.CreateFuel(fuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, fuel)
}

func (ctrl *FuelController) UpdateFuel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid fuel ID"})
	}

	var fuelRequest models.FuelRequest
	if err := c.Bind(&fuelRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if fuelRequest.Category == "" || fuelRequest.Name == "" || fuelRequest.EmissionFactor <= 0 || fuelRequest.Price <= 0 || fuelRequest.Unit == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	fuel, status, err := ctrl.FuelRepository.UpdateFuel(id, fuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, fuel)
}

func (ctrl *FuelController) DeleteFuel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid fuel ID"})
	}

	status, err := ctrl.FuelRepository.DeleteFuel(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete fuel"})
}
