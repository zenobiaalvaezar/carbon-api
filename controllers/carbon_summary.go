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

// GetCarbonSummary godoc
// @Summary Get carbon summary
// @Description Get the carbon summary for a specific user, including fuel and electric emissions, total emissions, and total trees equivalent
// @Tags CarbonSummary
// @Accept json
// @Produce json
// @Success 200 {object} models.CarbonSummaryResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /carbon-summary [get]
func (ctrl *CarbonSummaryController) GetCarbonSummary(c echo.Context) error {
	// userId := c.Get("user_id").(int)
	userId := 2 // Hardcoded user ID for testing

	carbonSummary, status, err := ctrl.CarbonSummaryRepository.GetCarbonSummary(userId)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, carbonSummary)
}
