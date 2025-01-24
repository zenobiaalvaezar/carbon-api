package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCart_Success(t *testing.T) {
	e := echo.New()
	mockRepository := new(repositories.MockCartRepository)

	mockRepository.On("GetAllCart", 1).Return([]models.GetCartsResponse{
		{ID: 1, TreeName: "Oak", Price: 25.5, Quantity: 2},
	}, http.StatusOK, nil)

	controller := NewCartController(mockRepository)

	req := httptest.NewRequest(http.MethodGet, "/carts", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("user_id", 1)

	err := controller.GetAllCart(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `[{"id":1,"tree_name":"Oak","price":25.5,"quantity":2}]`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}

func TestDeleteCart_Success(t *testing.T) {
	e := echo.New()
	mockRepository := new(repositories.MockCartRepository)

	mockRepository.On("DeleteCart", 1, 1).Return(http.StatusOK, nil)

	controller := NewCartController(mockRepository)

	req := httptest.NewRequest(http.MethodDelete, "/carts/1", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	ctx.Set("user_id", 1)

	err := controller.DeleteCart(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `{"message":"Success remove tree from cart"}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}
