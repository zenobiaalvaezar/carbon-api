package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
	RoleRepository repositories.RoleRepository
}

func NewRoleController(roleRepository repositories.RoleRepository) *RoleController {
	return &RoleController{roleRepository}
}

func (ctrl *RoleController) GetAllRoles(c echo.Context) error {
	roles, status, err := ctrl.RoleRepository.GetAllRoles()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, roles)
}

func (ctrl *RoleController) GetRoleByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	role, status, err := ctrl.RoleRepository.GetRoleByID(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, role)
}

func (ctrl *RoleController) CreateRole(c echo.Context) error {
	var roleRequest models.RoleRequest
	if err := c.Bind(&roleRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	role, status, err := ctrl.RoleRepository.CreateRole(roleRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, role)
}

func (ctrl *RoleController) UpdateRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	var roleRequest models.RoleRequest
	if err := c.Bind(&roleRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	role, status, err := ctrl.RoleRepository.UpdateRole(id, roleRequest)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, role)
}

func (ctrl *RoleController) DeleteRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role ID"})
	}

	status, err := ctrl.RoleRepository.DeleteRole(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, map[string]string{"message": "Role deleted successfully"})
}
