package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	PaymentRepository repositories.PaymentRepository
}

func NewPaymentController(paymentRepository repositories.PaymentRepository) *PaymentController {
	return &PaymentController{paymentRepository}
}

func (ctrl *PaymentController) CreatePayment(c echo.Context) error {
	var payment models.PaymentRequest
	c.Bind(&payment)

	userId := c.Get("user_id").(int)
	payment.UserID = userId

	paymentResponse, status, err := ctrl.PaymentRepository.CreatePayment(payment)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentResponse)
}

func (ctrl *PaymentController) VerifyPayment(c echo.Context) error {
	transactionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid transaction ID"})
	}

	status := c.QueryParam("status")
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Status is required"})
	}

	// if status is not success or failed
	if status != "success" && status != "failed" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid payment status"})
	}

	statusCode, err := ctrl.PaymentRepository.VerifyPayment(transactionID, status)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, map[string]string{"message": "Success verify payment"})
}
