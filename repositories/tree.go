package repositories

import (
	"carbon-api/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type TreeRepository interface {
	GetAllTrees() ([]models.Tree, int, error)
	GetTreeByID(id int) (models.Tree, int, error)
	CreateTree(tree *models.Tree) (int, error)
	UpdateTree(tree *models.Tree) (int, error)
	DeleteTree(id int) (int, error)
}

type treeRepository struct {
	DB *gorm.DB
}

func NewTreeRepository(DB *gorm.DB) TreeRepository {
	return &treeRepository{DB}
}

func (repo *treeRepository) GetAllTrees() ([]models.Tree, int, error) {
	var trees []models.Tree
	result := repo.DB.Order("id asc").Find(&trees)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}

	return trees, http.StatusOK, nil
}

func (repo *treeRepository) GetTreeByID(id int) (models.Tree, int, error) {
	var tree models.Tree
	result := repo.DB.First(&tree, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tree, http.StatusNotFound, errors.New("tree not found")
	}

	return tree, http.StatusOK, nil
}

func (repo *treeRepository) CreateTree(tree *models.Tree) (int, error) {
	// check if tree already exists
	var existingTree models.Tree
	result := repo.DB.Where("name = ?", tree.Name).First(&existingTree)
	if result.Error == nil {
		return http.StatusConflict, errors.New("tree already exists")
	}

	result = repo.DB.Create(tree)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusCreated, nil
}

func (repo *treeRepository) UpdateTree(tree *models.Tree) (int, error) {
	result := repo.DB.Save(tree)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (repo *treeRepository) DeleteTree(id int) (int, error) {
	result := repo.DB.Delete(&models.Tree{}, id)
	if result.RowsAffected == 0 {
		return http.StatusNotFound, errors.New("tree not found")
	}
	return http.StatusOK, nil
}
