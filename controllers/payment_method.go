package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"

	"github.com/labstack/echo/v4"
)

type PaymentMethodController struct {
	PaymentMethodRepository repositories.PaymentMethodRepository
}

func NewPaymentMethodController(paymentMethodRepository repositories.PaymentMethodRepository) *PaymentMethodController {
	return &PaymentMethodController{paymentMethodRepository}
}

func (ctrl *PaymentMethodController) GetAllPaymentMethods(c echo.Context) error {
	paymentMethods, status, err := ctrl.PaymentMethodRepository.GetAllPaymentMethods()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethods)
}

func (ctrl *PaymentMethodController) CreatePaymentMethod(c echo.Context) error {
	var paymentMethod models.PaymentMethod
	c.Bind(&paymentMethod)

	paymentMethodResponse, status, err := ctrl.PaymentMethodRepository.CreatePaymentMethod(paymentMethod)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethodResponse)
}

func (ctrl *PaymentMethodController) UpdatePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod
	c.Bind(&paymentMethod)

	paymentMethodResponse, status, err := ctrl.PaymentMethodRepository.UpdatePaymentMethod(id, paymentMethod)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, paymentMethodResponse)
}

func (ctrl *PaymentMethodController) DeletePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	status, err := ctrl.PaymentMethodRepository.DeletePaymentMethod(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	return c.JSON(status, map[string]string{"message": "Success delete payment method"})
}
