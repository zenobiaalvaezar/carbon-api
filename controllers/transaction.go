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

// GetAllTransactions godoc
// @Summary Get all transactions for a user
// @Description Retrieve a list of all transactions for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
func (controller *TransactionController) GetAllTransactions(c echo.Context) error {
	userId := c.Get("user_id").(int)
	transactions, statusCode, err := controller.TransactionRepository.GetAllTransactions(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, transactions)
}

// AddTransaction godoc
// @Summary Add a new transaction
// @Description Create a new transaction for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction data"
// @Success 200 {object} models.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /transactions [post]
func (controller *TransactionController) AddTransaction(c echo.Context) error {
	userId := c.Get("user_id").(int)
	transaction, statusCode, err := controller.TransactionRepository.AddTransaction(userId)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, transaction)
}
