package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ElectricController struct {
	ElectricRepository repositories.ElectricRepository
}

func NewElectricController(electricRepository repositories.ElectricRepository) *ElectricController {
	return &ElectricController{electricRepository}
}

func (ctrl *ElectricController) CreateElectric(c echo.Context) error {
	var electric models.Electric

	if err := c.Bind(&electric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	if err := ctrl.ElectricRepository.Create(&electric); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, electric)
}

func (ctrl *ElectricController) GetElectricByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	electric, err := ctrl.ElectricRepository.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if electric == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Record not found"})
	}

	return c.JSON(http.StatusOK, electric)
}

func (ctrl *ElectricController) GetAllElectrics(c echo.Context) error {
	electrics, err := ctrl.ElectricRepository.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, electrics)
}

func (ctrl *ElectricController) UpdateElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	var updatedElectric models.Electric
	if err := c.Bind(&updatedElectric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	if err := ctrl.ElectricRepository.Update(id, &updatedElectric); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Record updated successfully"})
}

func (ctrl *ElectricController) DeleteElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	if err := ctrl.ElectricRepository.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Record deleted successfully"})
}
