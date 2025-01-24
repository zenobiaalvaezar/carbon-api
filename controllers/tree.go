package controllers

import (
	"carbon-api/caches"
	"carbon-api/models"
	"carbon-api/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TreeController struct {
	TreeRepository repositories.TreeRepository
	TreeCache      caches.TreeCache
}

func NewTreeController(treeRepository repositories.TreeRepository, treeCache caches.TreeCache) *TreeController {
	return &TreeController{treeRepository, treeCache}
}

// GetAllTrees godoc
// @Summary Get all trees
// @Description Retrieve a list of all trees, checking consistency with the cache and database
// @Tags Trees
// @Accept json
// @Produce json
// @Success 200 {array} models.Tree
// @Failure 500 {object} map[string]string
// @Router /trees [get]
func (ctrl *TreeController) GetAllTrees(c echo.Context) error {
	var trees []models.Tree
	var err error

	// Step 1: Get all data from Redis
	trees, _, err = ctrl.TreeCache.GetAllTrees()
	if err == nil && len(trees) > 0 {
		// Step 2: Verify data consistency with the database
		dbTrees, _, dbErr := ctrl.TreeRepository.GetAllTrees()
		if dbErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": dbErr.Error()})
		}

		// Check if Redis and DB are in sync
		if len(trees) != len(dbTrees) {
			// If not in sync, update Redis
			_ = ctrl.TreeCache.SetAllTrees(dbTrees)
		}

		return c.JSON(http.StatusOK, dbTrees)
	}

	// Step 3: If Redis is empty, get all data from database
	trees, status, err := ctrl.TreeRepository.GetAllTrees()
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Step 4: Sync all data to Redis
	_ = ctrl.TreeCache.SetAllTrees(trees)

	return c.JSON(http.StatusOK, trees)
}

// GetTreeByID godoc
// @Summary Get tree by ID
// @Description Retrieve a specific tree by ID, checking cache first, then the database
// @Tags Trees
// @Accept json
// @Produce json
// @Param id path int true "Tree ID"
// @Success 200 {object} models.Tree
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /trees/{id} [get]
func (ctrl *TreeController) GetTreeByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Get from Redis
	tree, _, err := ctrl.TreeCache.GetTreeByID(id)
	if err == nil {
		return c.JSON(http.StatusOK, tree)
	}

	// Get from Database
	tree, status, err := ctrl.TreeRepository.GetTreeByID(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Cache the data
	_ = ctrl.TreeCache.SetTree(tree)

	return c.JSON(http.StatusOK, tree)
}

// CreateTree godoc
// @Summary Create a new tree
// @Description Add a new tree to the database and cache
// @Tags Trees
// @Accept json
// @Produce json
// @Param tree body models.Tree true "Tree data"
// @Success 201 {object} models.Tree
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trees [post]
func (ctrl *TreeController) CreateTree(c echo.Context) error {
	var tree models.Tree
	if err := c.Bind(&tree); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if tree.TreeCategoryID == 0 || tree.Name == "" || tree.Description == "" || tree.Price <= 0 || tree.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	status, err := ctrl.TreeRepository.CreateTree(&tree)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Cache the new tree
	_ = ctrl.TreeCache.SetTree(tree)

	return c.JSON(http.StatusCreated, tree)
}

// UpdateTree godoc
// @Summary Update an existing tree
// @Description Modify an existing tree by ID and update the cache
// @Tags Trees
// @Accept json
// @Produce json
// @Param id path int true "Tree ID"
// @Param tree body models.Tree true "Tree data"
// @Success 200 {object} models.Tree
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trees/{id} [put]
func (ctrl *TreeController) UpdateTree(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var tree models.Tree
	if err := c.Bind(&tree); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if tree.TreeCategoryID == 0 || tree.Name == "" || tree.Description == "" || tree.Price <= 0 || tree.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	tree.ID = id
	status, err := ctrl.TreeRepository.UpdateTree(&tree)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Update cache
	_ = ctrl.TreeCache.SetTree(tree)

	return c.JSON(http.StatusOK, tree)
}

// DeleteTree godoc
// @Summary Delete a tree
// @Description Remove a tree by ID from the database and cache
// @Tags Trees
// @Accept json
// @Produce json
// @Param id path int true "Tree ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /trees/{id} [delete]
func (ctrl *TreeController) DeleteTree(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Delete from Database
	status, err := ctrl.TreeRepository.DeleteTree(id)
	if err != nil {
		return c.JSON(status, map[string]string{"message": err.Error()})
	}

	// Delete from Redis
	_ = ctrl.TreeCache.DeleteTree(id)

	return c.JSON(http.StatusOK, map[string]string{"message": "Tree deleted"})
}
