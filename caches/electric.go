package caches

import (
	"carbon-api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
)

var electricKey = "electrics"

type ElectricCache interface {
	GetAllElectrics() ([]models.Electric, int, error)
	GetElectricByID(id int) (models.Electric, int, error)
	CreateAllElectrics(electrics []models.Electric) (int, error)
	CreateElectric(electric models.Electric) (int, error)
	UpdateElectric(electric models.Electric) (int, error)
	DeleteElectric(id int) (int, error)
}

type electricCache struct {
	client *redis.Client
}

func NewElectricCache(client *redis.Client) ElectricCache {
	return &electricCache{client}
}

func serializeElectric(electric models.Electric) (string, error) {
	data, err := json.Marshal(electric)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func deserializeElectric(data string) (models.Electric, error) {
	var electric models.Electric
	if err := json.Unmarshal([]byte(data), &electric); err != nil {
		return models.Electric{}, err
	}
	return electric, nil
}

func (cache *electricCache) GetAllElectrics() ([]models.Electric, int, error) {
	vals, err := cache.client.HGetAll(electricKey).Result()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	electrics := make([]models.Electric, 0, len(vals))
	for _, val := range vals {
		electric, err := deserializeElectric(val)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		electrics = append(electrics, electric)
	}

	return electrics, http.StatusOK, nil
}

func (cache *electricCache) GetElectricByID(id int) (models.Electric, int, error) {
	val, err := cache.client.HGet(electricKey, strconv.Itoa(id)).Result()
	if err != nil {
		return models.Electric{}, http.StatusNotFound, errors.New("electric record not found")
	}

	electric, err := deserializeElectric(val)
	if err != nil {
		return models.Electric{}, http.StatusInternalServerError, err
	}

	return electric, http.StatusOK, nil
}

func (cache *electricCache) CreateAllElectrics(electrics []models.Electric) (int, error) {
	pipe := cache.client.Pipeline()

	for _, electric := range electrics {
		data, err := serializeElectric(electric)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		pipe.HSet(electricKey, strconv.Itoa(electric.ID), data)
	}

	_, err := pipe.Exec()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (cache *electricCache) CreateElectric(electric models.Electric) (int, error) {
	data, err := serializeElectric(electric)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cache.client.HSet(electricKey, strconv.Itoa(electric.ID), data).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (cache *electricCache) UpdateElectric(electric models.Electric) (int, error) {
	_, err := cache.client.HExists(electricKey, strconv.Itoa(electric.ID)).Result()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data, err := serializeElectric(electric)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cache.client.HSet(electricKey, strconv.Itoa(electric.ID), data).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cache *electricCache) DeleteElectric(id int) (int, error) {
	deleted, err := cache.client.HDel(electricKey, strconv.Itoa(id)).Result()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if deleted == 0 {
		return http.StatusNotFound, errors.New("electric record not found")
	}

	return http.StatusOK, nil
}
