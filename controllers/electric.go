package controllers

import (
	"carbon-api/caches"
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ElectricController struct {
	ElectricRepository repositories.ElectricRepository
	ElectricCache      caches.ElectricCache
}

func NewElectricController(electricRepository repositories.ElectricRepository, electricCache caches.ElectricCache) *ElectricController {
	return &ElectricController{electricRepository, electricCache}
}

// GetAllElectrics godoc
// @Summary Get all electric data
// @Description Retrieve a list of all electric records
// @Tags Electrics
// @Accept json
// @Produce json
// @Success 200 {array} models.Electric "List of electrics"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /electrics [get]
func (ctrl *ElectricController) GetAllElectrics(c echo.Context) error {
	electrics, status, err := ctrl.ElectricCache.GetAllElectrics()
	if err != nil || len(electrics) == 0 {
		electrics, err = ctrl.ElectricRepository.FindAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		cacheStatus, cacheErr := ctrl.ElectricCache.CreateAllElectrics(electrics)
		if cacheErr != nil {
			return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
		}
	}
	return c.JSON(status, electrics)
}

// GetElectricByID godoc
// @Summary Get an electric record by ID
// @Description Fetch an electric record using its unique ID
// @Tags Electrics
// @Accept json
// @Produce json
// @Param id path int true "Electric ID"
// @Success 200 {object} models.Electric "Electric details"
// @Failure 400 {object} map[string]interface{} "Invalid electric ID"
// @Failure 404 {object} map[string]interface{} "Electric not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /electrics/{id} [get]
func (ctrl *ElectricController) GetElectricByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid electric ID"})
	}

	electric, err := ctrl.ElectricRepository.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if electric == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Electric not found"})
	}

	electricCopy := *electric

	cacheStatus, cacheErr := ctrl.ElectricCache.CreateElectric(electricCopy)
	if cacheErr != nil {
		return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
	}

	return c.JSON(cacheStatus, electric)
}

// CreateElectric godoc
// @Summary Create a new electric record
// @Description Add a new electric record to the database
// @Tags Electrics
// @Accept json
// @Produce json
// @Param electric body models.Electric true "Electric data"
// @Success 201 {object} models.Electric "Created electric"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /electrics [post]
func (ctrl *ElectricController) CreateElectric(c echo.Context) error {
	var electric models.Electric
	if err := c.Bind(&electric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	err := ctrl.ElectricRepository.Create(&electric)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	cacheStatus, cacheErr := ctrl.ElectricCache.CreateElectric(electric)
	if cacheErr != nil {
		return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
	}

	return c.JSON(http.StatusCreated, electric)
}

// UpdateElectric godoc
// @Summary Update an electric record by ID
// @Description Modify the details of an existing electric record
// @Tags Electrics
// @Accept json
// @Produce json
// @Param id path int true "Electric ID"
// @Param electric body models.Electric true "Updated electric data"
// @Success 200 {object} models.Electric "Updated electric"
// @Failure 400 {object} map[string]interface{} "Invalid electric ID or payload"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /electrics/{id} [put]
func (ctrl *ElectricController) UpdateElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid electric ID"})
	}

	var electric models.Electric
	if err := c.Bind(&electric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	err = ctrl.ElectricRepository.Update(id, &electric)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	electric.ID = id
	cacheStatus, cacheErr := ctrl.ElectricCache.UpdateElectric(electric)
	if cacheErr != nil {
		return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
	}

	return c.JSON(http.StatusOK, electric)
}

// DeleteElectric godoc
// @Summary Delete an electric record by ID
// @Description Remove an electric record from the database
// @Tags Electrics
// @Accept json
// @Produce json
// @Param id path int true "Electric ID"
// @Success 200 {object} map[string]interface{} "Electric deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid electric ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /electrics/{id} [delete]
func (ctrl *ElectricController) DeleteElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid electric ID"})
	}

	err = ctrl.ElectricRepository.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	cacheStatus, cacheErr := ctrl.ElectricCache.DeleteElectric(id)
	if cacheErr != nil {
		return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Electric deleted successfully"})
}
