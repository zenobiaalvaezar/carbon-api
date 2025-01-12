package controllers

import (
	"carbon-api/caches"
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FuelController struct {
	FuelRepository repositories.FuelRepository
	FuelCache      caches.FuelCache
}

func NewFuelController(fuelRepository repositories.FuelRepository, fuelCache caches.FuelCache) *FuelController {
	return &FuelController{fuelRepository, fuelCache}
}

func (ctrl *FuelController) GetAllFuels(c echo.Context) error {
	var fuels []models.Fuel
	var status int

	// Get all fuels from cache
	fuels, status, err := ctrl.FuelCache.GetAllFuels()

	if err != nil || len(fuels) == 0 {
		// Get all fuels from database
		fuels, status, err = ctrl.FuelRepository.GetAllFuels()

		if err != nil {
			return c.JSON(status, map[string]string{"message": err.Error()})
		}

		// Save all fuels to cache
		cacheStatus, err := ctrl.FuelCache.CreateAllFuels(fuels)
		if err != nil {
			return c.JSON(cacheStatus, map[string]string{"message": err.Error()})
		}
	}

	return c.JSON(status, fuels)
}

func (ctrl *FuelController) GetFuelByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid fuel ID"})
	}

	var fuel models.Fuel
	var status int

	// Get fuel from cache
	fuel, status, err = ctrl.FuelCache.GetFuelByID(id)

	if err != nil || fuel.ID == 0 {
		// Get fuel from database
		fuel, status, err = ctrl.FuelRepository.GetFuelByID(id)

		if err != nil {
			return c.JSON(status, map[string]string{"message": err.Error()})
		}

		// Save fuel to cache
		cacheStatus, err := ctrl.FuelCache.CreateFuel(fuel)
		if err != nil {
			return c.JSON(cacheStatus, map[string]string{"message": err.Error()})
		}
	}

	return c.JSON(status, fuel)
}

func (ctrl *FuelController) CreateFuel(c echo.Context) error {
	var fuelRequest models.FuelRequest
	if err := c.Bind(&fuelRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	if fuelRequest.Category == "" || fuelRequest.Name == "" || fuelRequest.EmissionFactor <= 0 || fuelRequest.Price <= 0 || fuelRequest.Unit == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Create fuel in database
	fuel, status, err := ctrl.FuelRepository.CreateFuel(fuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Save fuel to cache
	cacheStatus, err := ctrl.FuelCache.CreateFuel(fuel)
	if err != nil {
		return c.JSON(cacheStatus, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, fuel)
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

	// Update fuel in database
	fuel, status, err := ctrl.FuelRepository.UpdateFuel(id, fuelRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Save fuel to cache
	cacheStatus, err := ctrl.FuelCache.UpdateFuel(fuel)
	if err != nil {
		return c.JSON(cacheStatus, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, fuel)
}

func (ctrl *FuelController) DeleteFuel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid fuel ID"})
	}

	// Delete fuel in database
	status, err := ctrl.FuelRepository.DeleteFuel(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Delete fuel from cache
	cacheStatus, err := ctrl.FuelCache.DeleteFuel(id)
	if err != nil {
		return c.JSON(cacheStatus, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete fuel"})
}
