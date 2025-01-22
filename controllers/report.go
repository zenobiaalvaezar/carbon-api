package controllers

import (
	"carbon-api/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReportController struct {
	ReportRepository repositories.ReportRepository
}

func NewReportController(reportRepository repositories.ReportRepository) *ReportController {
	return &ReportController{reportRepository}
}

// GetReportSummary godoc
// @Summary Get the report summary for a user
// @Description Retrieve a summary report for the authenticated user, including user details and statistics such as carbon trees and donation trees
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} models.ReportSummary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /reports/summary [get]
func (controller *ReportController) GetReportSummary(c echo.Context) error {
	userId := c.Get("user_id").(int)
	report, statusCode, err := controller.ReportRepository.GetReportSummary(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, report)
}
