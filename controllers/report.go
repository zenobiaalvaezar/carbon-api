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

func (controller *ReportController) GetReportSummary(c echo.Context) error {
	userId := c.Get("user_id").(int)
	report, statusCode, err := controller.ReportRepository.GetReportSummary(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, report)
}
