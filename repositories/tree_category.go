package repositories

import (
	"carbon-api/config"
	"carbon-api/models"
	"errors"
	"strings"
)

func GetAllTreeCategories() ([]models.TreeCategory, error) {
	var categories []models.TreeCategory
	err := config.DB.Find(&categories).Error
	return categories, err
}

func GetTreeCategoryByID(id int) (models.TreeCategory, error) {
	var category models.TreeCategory
	err := config.DB.First(&category, id).Error
	return category, err
}

func CreateTreeCategory(category *models.TreeCategory) error {
	err := config.DB.Create(category).Error
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return errors.New("category name already exists")
	}
	return err
}

func UpdateTreeCategory(category *models.TreeCategory) error {
	// Cek apakah ID ada di database
	var existingCategory models.TreeCategory
	if err := config.DB.First(&existingCategory, category.ID).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return errors.New("category with the given ID does not exist")
		}
		return err
	}

	// Lanjutkan update jika ID ditemukan
	return config.DB.Save(category).Error
}

func DeleteTreeCategory(id int) error {
	return config.DB.Delete(&models.TreeCategory{}, id).Error
}
