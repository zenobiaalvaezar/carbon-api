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

// GetAllRoles godoc
// @Summary Get all roles
// @Description Retrieve a list of all roles
// @Tags Roles
// @Accept json
// @Produce json
// @Success 200 {array} models.Role "List of roles"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles [get]
func (ctrl *RoleController) GetAllRoles(c echo.Context) error {
	roles, status, err := ctrl.RoleRepository.GetAllRoles()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}
	return c.JSON(status, roles)
}

// GetRoleByID godoc
// @Summary Get a role by ID
// @Description Retrieve details of a specific role by its unique ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role "Role details"
// @Failure 400 {object} map[string]string "Invalid role ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [get]
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

// CreateRole godoc
// @Summary Create a new role
// @Description Add a new role to the system
// @Tags Roles
// @Accept json
// @Produce json
// @Param roleRequest body models.RoleRequest true "Role request payload"
// @Success 201 {object} models.Role "Successfully created role"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles [post]
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

// UpdateRole godoc
// @Summary Update a role
// @Description Update details of an existing role by its unique ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param roleRequest body models.RoleRequest true "Role update request payload"
// @Success 200 {object} models.Role "Successfully updated role"
// @Failure 400 {object} map[string]string "Invalid role ID or request payload"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [put]
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

// DeleteRole godoc
// @Summary Delete a role
// @Description Remove a specific role by its unique ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Successfully deleted role"
// @Failure 400 {object} map[string]string "Invalid role ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /roles/{id} [delete]
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
