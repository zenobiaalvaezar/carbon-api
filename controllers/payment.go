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

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a new payment request for a user, specifying the transaction ID, payment method, and other details
// @Tags Payments
// @Accept json
// @Produce json
// @Param paymentRequest body models.PaymentRequest true "Payment Request"
// @Success 200 {object} models.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payments [post]
func (ctrl *PaymentController) CreatePayment(c echo.Context) error {
	var payment models.PaymentRequest
	c.Bind(&payment)

	if payment.TransactionID == 0 || payment.PaymentMethod == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	userId := c.Get("user_id").(int)
	payment.UserID = userId

	paymentResponse, status, err := ctrl.PaymentRepository.CreatePayment(payment)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentResponse)
}

// VerifyPayment godoc
// @Summary Verify a payment transaction
// @Description Verify the status of a payment transaction by its transaction ID and status (either success or failed)
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param status query string true "Payment Status" Enums(success, failed)
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payments/verify/{id} [get]
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
