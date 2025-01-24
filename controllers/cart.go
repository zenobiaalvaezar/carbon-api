package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartController struct {
	CartRepository repositories.CartRepository
}

func NewCartController(cr repositories.CartRepository) *CartController {
	return &CartController{cr}
}

// GetAllCart godoc
// @Summary Get all carts for a specific user
// @Description Retrieve all cart items for a specific user by their user ID
// @Tags Carts
// @Accept json
// @Produce json
// @Success 200 {array} models.GetCartsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carts [get]
func (cc *CartController) GetAllCart(c echo.Context) error {
	userId := c.Get("user_id").(int)
	carts, statusCode, err := cc.CartRepository.GetAllCart(userId)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, carts)
}

// AddCart godoc
// @Summary Add a tree to the cart
// @Description Add a specified tree item to the user's cart
// @Tags Carts
// @Accept json
// @Produce json
// @Param addCartRequest body models.AddCartRequest true "Add Tree to Cart Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carts [post]
func (cc *CartController) AddCart(c echo.Context) error {
	var cart models.AddCartRequest
	c.Bind(&cart)

	if cart.TreeID == 0 || cart.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	userId := c.Get("user_id").(int)
	cart.UserID = userId

	newCart, statusCode, err := cc.CartRepository.AddCart(cart)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, newCart)
}

// DeleteCart godoc
// @Summary Delete a tree from the cart
// @Description Remove a specified tree item from the user's cart by ID
// @Tags Carts
// @Accept json
// @Produce json
// @Param id path int true "Cart ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carts/{id} [delete]
func (cc *CartController) DeleteCart(c echo.Context) error {
	id := c.Param("id")
	cartId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	userId := c.Get("user_id").(int)
	statusCode, err := cc.CartRepository.DeleteCart(cartId, userId)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, map[string]string{"message": "Success remove tree from cart"})
}
