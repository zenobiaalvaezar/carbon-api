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

func TestRoleController_GetAllRoles(t *testing.T) {
	mockRepo := new(repositories.MockRoleRepository)
	ctrl := NewRoleController(mockRepo)

	expectedRoles := []models.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "User"},
	}
	mockRepo.On("GetAllRoles").Return(expectedRoles, http.StatusOK, nil)

	req := httptest.NewRequest(http.MethodGet, "/roles", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := ctrl.GetAllRoles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"id":1,"name":"Admin"},{"id":2,"name":"User"}]`, rec.Body.String())

	mockRepo.AssertExpectations(t)
}
