package controllers

import (
	"carbon-api/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionRepository repositories.TransactionRepository
}

func NewTransactionController(TransactionRepository repositories.TransactionRepository) *TransactionController {
	return &TransactionController{TransactionRepository}
}

func (controller *TransactionController) GetAllTransactions(c echo.Context) error {
	userId := c.Get("user").(int)
	transactions, statusCode, err := controller.TransactionRepository.GetAllTransactions(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, transactions)
}

func (controller *TransactionController) AddTransaction(c echo.Context) error {
	userId := c.Get("user").(int)
	transaction, statusCode, err := controller.TransactionRepository.AddTransaction(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, transaction)
}
