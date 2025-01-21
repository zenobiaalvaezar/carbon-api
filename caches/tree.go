package caches

import (
	"carbon-api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
)

var treeKey = "trees"

type TreeCache interface {
	GetAllTrees() ([]models.Tree, int, error)
	GetTreeByID(id int) (models.Tree, int, error)
	SetAllTrees(trees []models.Tree) error
	SetTree(tree models.Tree) error
	DeleteTree(id int) error
}

type treeCache struct {
	client *redis.Client
}

func NewTreeCache(client *redis.Client) TreeCache {
	return &treeCache{client}
}

func serializeTree(tree models.Tree) (string, error) {
	data, err := json.Marshal(tree)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func deserializeTree(data string) (models.Tree, error) {
	var tree models.Tree
	if err := json.Unmarshal([]byte(data), &tree); err != nil {
		return models.Tree{}, err
	}
	return tree, nil
}

func (cache *treeCache) GetAllTrees() ([]models.Tree, int, error) {
	vals, err := cache.client.HGetAll(treeKey).Result()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	trees := make([]models.Tree, 0, len(vals))
	for _, val := range vals {
		tree, err := deserializeTree(val)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		trees = append(trees, tree)
	}

	return trees, http.StatusOK, nil
}

func (cache *treeCache) GetTreeByID(id int) (models.Tree, int, error) {
	val, err := cache.client.HGet(treeKey, strconv.Itoa(id)).Result()
	if err != nil {
		return models.Tree{}, http.StatusNotFound, errors.New("tree not found in cache")
	}

	tree, err := deserializeTree(val)
	if err != nil {
		return models.Tree{}, http.StatusInternalServerError, err
	}

	return tree, http.StatusOK, nil
}

func (cache *treeCache) SetAllTrees(trees []models.Tree) error {
	pipe := cache.client.Pipeline()

	for _, tree := range trees {
		data, err := serializeTree(tree)
		if err != nil {
			return err
		}
		pipe.HSet(treeKey, strconv.Itoa(tree.ID), data)
	}

	_, err := pipe.Exec()
	return err
}

func (cache *treeCache) SetTree(tree models.Tree) error {
	data, err := serializeTree(tree)
	if err != nil {
		return err
	}
	return cache.client.HSet(treeKey, strconv.Itoa(tree.ID), data).Err()
}

func (cache *treeCache) DeleteTree(id int) error {
	return cache.client.HDel(treeKey, strconv.Itoa(id)).Err()
}
