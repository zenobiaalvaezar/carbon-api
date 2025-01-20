package controllers

import (
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
