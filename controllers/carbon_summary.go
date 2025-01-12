package controllers

import (
	"carbon-api/repositories"

	"github.com/labstack/echo/v4"
)

type CarbonSummaryController struct {
	CarbonSummaryRepository repositories.CarbonSummaryRepository
}

func NewCarbonSummaryController(carbonSummaryRepository repositories.CarbonSummaryRepository) *CarbonSummaryController {
	return &CarbonSummaryController{carbonSummaryRepository}
}

func (ctrl *CarbonSummaryController) GetCarbonSummary(c echo.Context) error {
	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing

	carbonSummary, status, err := ctrl.CarbonSummaryRepository.GetCarbonSummary(userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonSummary)
}
