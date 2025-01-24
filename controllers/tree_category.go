package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllTreeCategories godoc
// @Summary Get all tree categories
// @Description Retrieve a list of all tree categories
// @Tags TreeCategories
// @Accept json
// @Produce json
// @Success 200 {array} models.TreeCategory
// @Failure 500 {object} map[string]string
// @Router /tree-categories [get]
func GetAllTreeCategories(c echo.Context) error {
	categories, err := repositories.GetAllTreeCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}

// GetTreeCategoryByID godoc
// @Summary Get tree category by ID
// @Description Retrieve a tree category by its ID
// @Tags TreeCategories
// @Accept json
// @Produce json
// @Param id path int true "Tree Category ID"
// @Success 200 {object} models.TreeCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tree-categories/{id} [get]
func GetTreeCategoryByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := repositories.GetTreeCategoryByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	return c.JSON(http.StatusOK, category)
}

// CreateTreeCategory godoc
// @Summary Create a new tree category
// @Description Create a new tree category with a unique name
// @Tags TreeCategories
// @Accept json
// @Produce json
// @Param category body models.TreeCategory true "Tree Category"
// @Success 201 {object} models.TreeCategory
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tree-categories [post]
func CreateTreeCategory(c echo.Context) error {
	var category models.TreeCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if category.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
	}

	err := repositories.CreateTreeCategory(&category)
	if err != nil {
		if err.Error() == "category name already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Category name already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, category)
}

// UpdateTreeCategory godoc
// @Summary Update an existing tree category
// @Description Update a tree category by ID
// @Tags TreeCategories
// @Accept json
// @Produce json
// @Param id path int true "Tree Category ID"
// @Param category body models.TreeCategory true "Tree Category"
// @Success 200 {object} models.TreeCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tree-categories/{id} [put]
func UpdateTreeCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.TreeCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if category.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
	}

	category.ID = id
	err := repositories.UpdateTreeCategory(&category)
	if err != nil {
		if err.Error() == "category with the given ID does not exist" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category with the given ID does not exist"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// DeleteTreeCategory godoc
// @Summary Delete a tree category
// @Description Delete a tree category by ID
// @Tags TreeCategories
// @Accept json
// @Produce json
// @Param id path int true "Tree Category ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tree-categories/{id} [delete]
func DeleteTreeCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repositories.DeleteTreeCategory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
