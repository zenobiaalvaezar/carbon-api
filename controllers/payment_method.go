package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentMethodController struct {
	PaymentMethodRepository repositories.PaymentMethodRepository
}

func NewPaymentMethodController(paymentMethodRepository repositories.PaymentMethodRepository) *PaymentMethodController {
	return &PaymentMethodController{paymentMethodRepository}
}

// GetAllPaymentMethods godoc
// @Summary Get all payment methods
// @Description Fetch all payment methods from the database
// @Tags PaymentMethods
// @Accept json
// @Produce json
// @Success 200 {array} models.PaymentMethod
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payment-methods [get]
func (ctrl *PaymentMethodController) GetAllPaymentMethods(c echo.Context) error {
	paymentMethods, status, err := ctrl.PaymentMethodRepository.GetAllPaymentMethods()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethods)
}

// CreatePaymentMethod godoc
// @Summary Create a new payment method
// @Description Create a new payment method by providing the necessary details
// @Tags PaymentMethods
// @Accept json
// @Produce json
// @Param paymentMethod body models.PaymentMethod true "New Payment Method"
// @Success 201 {object} models.PaymentMethod
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payment-methods [post]
func (ctrl *PaymentMethodController) CreatePaymentMethod(c echo.Context) error {
	var paymentMethod models.PaymentMethod
	c.Bind(&paymentMethod)

	if paymentMethod.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Code is required"})
	}

	paymentMethodResponse, status, err := ctrl.PaymentMethodRepository.CreatePaymentMethod(paymentMethod)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethodResponse)
}

// UpdatePaymentMethod godoc
// @Summary Update an existing payment method
// @Description Update the details of an existing payment method by its ID
// @Tags PaymentMethods
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Param paymentMethod body models.PaymentMethod true "Updated Payment Method"
// @Success 200 {object} models.PaymentMethod
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payment-methods/{id} [put]
func (ctrl *PaymentMethodController) UpdatePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod
	c.Bind(&paymentMethod)

	if paymentMethod.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Code is required"})
	}

	paymentMethodResponse, status, err := ctrl.PaymentMethodRepository.UpdatePaymentMethod(id, paymentMethod)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethodResponse)
}

// DeletePaymentMethod godoc
// @Summary Delete a payment method
// @Description Delete an existing payment method by its ID
// @Tags PaymentMethods
// @Accept json
// @Produce json
// @Param id path string true "Payment Method ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /payment-methods/{id} [delete]
func (ctrl *PaymentMethodController) DeletePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	status, err := ctrl.PaymentMethodRepository.DeletePaymentMethod(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete payment method"})
}
