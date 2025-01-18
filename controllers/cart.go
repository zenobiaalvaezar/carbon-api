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

func (cc *CartController) GetAllCart(c echo.Context) error {
	userId := c.Get("user_id").(int)
	carts, statusCode, err := cc.CartRepository.GetAllCart(userId)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, carts)
}

func (cc *CartController) AddCart(c echo.Context) error {
	var cart models.AddCartRequest
	c.Bind(&cart)

	userId := c.Get("user_id").(int)
	cart.UserID = userId

	statusCode, err := cc.CartRepository.AddCart(cart)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, map[string]string{"message": "Success add tree to cart"})
}

func (cc *CartController) DeleteCart(c echo.Context) error {
	id := c.Param("id")
	cartId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	statusCode, err := cc.CartRepository.DeleteCart(cartId)
	if err != nil {
		return c.JSON(statusCode, map[string]string{"message": err.Error()})
	}

	return c.JSON(statusCode, map[string]string{"message": "Success remove tree from cart"})
}
