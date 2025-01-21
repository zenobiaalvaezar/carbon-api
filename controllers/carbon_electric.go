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

// GetAllCarbonElectrics godoc
// @Summary Get all carbon electrics
// @Description Retrieve a list of carbon electric entries for the logged-in user
// @Tags CarbonElectrics
// @Accept json
// @Produce json
// @Success 200 {array} models.CarbonElectricResponse "List of carbon electric entries"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /carbon-electrics [get]
func (ctrl *CarbonElectricController) GetAllCarbonElectrics(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User ID is invalid or missing. Please log in to continue.",
		})
	}

	carbonElectrics, status, err := ctrl.CarbonElectricRepository.GetAllCarbonElectrics(userID)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonElectrics)
}

// GetCarbonElectricByID godoc
// @Summary Get a carbon electric entry by ID
// @Description Fetch details of a specific carbon electric entry by its unique ID
// @Tags CarbonElectrics
// @Accept json
// @Produce json
// @Param id path int true "Carbon Electric ID"
// @Success 200 {object} models.CarbonElectricResponse "Carbon electric entry details"
// @Failure 400 {object} map[string]string "Invalid carbon electric ID"
// @Failure 404 {object} map[string]string "Carbon electric not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /carbon-electrics/{id} [get]
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

// CreateCarbonElectric godoc
// @Summary Create a new carbon electric entry
// @Description Add a new carbon electric entry for the logged-in user
// @Tags CarbonElectrics
// @Accept json
// @Produce json
// @Param carbonElectricRequest body models.CarbonElectricRequest true "Carbon Electric request body"
// @Success 201 {object} models.CarbonElectricResponse "Successfully created carbon electric entry"
// @Failure 400 {object} map[string]string "Invalid request payload or usage type/amount"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /carbon-electrics [post]
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

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User ID is invalid or missing. Please log in to continue.",
		})
	}

	carbonElectricRequest.UserID = userID
	carbonElectric, status, err := ctrl.CarbonElectricRepository.CreateCarbonElectric(carbonElectricRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonElectric)
}

// DeleteCarbonElectric godoc
// @Summary Delete a carbon electric entry
// @Description Delete a specific carbon electric entry by its unique ID
// @Tags CarbonElectrics
// @Accept json
// @Produce json
// @Param id path int true "Carbon Electric ID"
// @Success 200 {object} map[string]string "Successfully deleted carbon electric entry"
// @Failure 400 {object} map[string]string "Invalid carbon electric ID"
// @Failure 404 {object} map[string]string "Carbon electric not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /carbon-electrics/{id} [delete]
func (ctrl *CarbonElectricController) DeleteCarbonElectric(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid carbon electric ID"})
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User ID is invalid or missing. Please log in to continue.",
		})
	}

	status, err := ctrl.CarbonElectricRepository.DeleteCarbonElectric(id, userID)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete carbon electric"})
}
