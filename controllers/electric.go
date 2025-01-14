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

	cacheStatus, cacheErr := ctrl.ElectricCache.UpdateElectric(electric)
	if cacheErr != nil {
		return c.JSON(cacheStatus, map[string]string{"message": cacheErr.Error()})
	}

	return c.JSON(http.StatusOK, electric)
}

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
