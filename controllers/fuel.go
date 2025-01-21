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

// GetAllFuels godoc
// @Summary Get all fuels
// @Description Retrieve a list of all fuels from cache, if not found, retrieve from database and cache it
// @Tags Fuels
// @Accept json
// @Produce json
// @Success 200 {array} models.Fuel
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /fuels [get]
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

// GetFuelByID godoc
// @Summary Get a fuel by ID
// @Description Retrieve a specific fuel by ID from cache, if not found, retrieve from database and cache it
// @Tags Fuels
// @Accept json
// @Produce json
// @Param id path int true "Fuel ID"
// @Success 200 {object} models.Fuel
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /fuels/{id} [get]
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

// CreateFuel godoc
// @Summary Create a new fuel
// @Description Create a new fuel record and store it in both the database and cache
// @Tags Fuels
// @Accept json
// @Produce json
// @Param body body models.FuelRequest true "Fuel Request"
// @Success 201 {object} models.Fuel
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /fuels [post]
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

// UpdateFuel godoc
// @Summary Update an existing fuel
// @Description Update a specific fuel by ID in both the database and cache
// @Tags Fuels
// @Accept json
// @Produce json
// @Param id path int true "Fuel ID"
// @Param body body models.FuelRequest true "Fuel Request"
// @Success 200 {object} models.Fuel
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /fuels/{id} [put]
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

// DeleteFuel godoc
// @Summary Delete a fuel
// @Description Delete a specific fuel by ID from both the database and cache
// @Tags Fuels
// @Accept json
// @Produce json
// @Param id path int true "Fuel ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /fuels/{id} [delete]
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
