package controllers

import (
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllTreeCategories(c echo.Context) error {
	categories, err := repositories.GetAllTreeCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}

func GetTreeCategoryByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := repositories.GetTreeCategoryByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	return c.JSON(http.StatusOK, category)
}

func CreateTreeCategory(c echo.Context) error {
	var category models.TreeCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
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

func UpdateTreeCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.TreeCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
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

func DeleteTreeCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repositories.DeleteTreeCategory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
