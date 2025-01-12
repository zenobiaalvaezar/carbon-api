package caches

import (
	"carbon-api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
)

var key = "fuels"

type FuelCache interface {
	GetAllFuels() ([]models.Fuel, int, error)
	GetFuelByID(id int) (models.Fuel, int, error)
	CreateAllFuels(fuels []models.Fuel) (int, error)
	CreateFuel(fuel models.Fuel) (int, error)
	UpdateFuel(fuel models.Fuel) (int, error)
	DeleteFuel(id int) (int, error)
}

type fuelCache struct {
	client *redis.Client
}

func NewFuelCache(client *redis.Client) FuelCache {
	return &fuelCache{client}
}

// serializeFuel converts a Fuel struct to JSON bytes
func serializeFuel(fuel models.Fuel) (string, error) {
	data, err := json.Marshal(fuel)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// deserializeFuel converts JSON bytes to a Fuel struct
func deserializeFuel(data string) (models.Fuel, error) {
	var fuel models.Fuel
	if err := json.Unmarshal([]byte(data), &fuel); err != nil {
		return models.Fuel{}, err
	}
	return fuel, nil
}

func (cache *fuelCache) GetAllFuels() ([]models.Fuel, int, error) {
	vals, err := cache.client.HGetAll(key).Result()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	fuels := make([]models.Fuel, 0, len(vals))
	for _, val := range vals {
		fuel, err := deserializeFuel(val)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		fuels = append(fuels, fuel)
	}

	return fuels, http.StatusOK, nil
}

func (cache *fuelCache) GetFuelByID(id int) (models.Fuel, int, error) {
	val, err := cache.client.HGet(key, strconv.Itoa(id)).Result()
	if err != nil {
		return models.Fuel{}, http.StatusNotFound, errors.New("fuel not found")
	}

	fuel, err := deserializeFuel(val)
	if err != nil {
		return models.Fuel{}, http.StatusInternalServerError, err
	}

	return fuel, http.StatusOK, nil
}

func (cache *fuelCache) CreateAllFuels(fuels []models.Fuel) (int, error) {
	pipe := cache.client.Pipeline()

	for _, fuel := range fuels {
		data, err := serializeFuel(fuel)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		pipe.HSet(key, strconv.Itoa(fuel.ID), data)
	}

	_, err := pipe.Exec()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (cache *fuelCache) CreateFuel(fuel models.Fuel) (int, error) {
	data, err := serializeFuel(fuel)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cache.client.HSet(key, strconv.Itoa(fuel.ID), data).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (cache *fuelCache) UpdateFuel(fuel models.Fuel) (int, error) {
	// Check if fuel exists
	exists, err := cache.client.HExists(key, strconv.Itoa(fuel.ID)).Result()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exists {
		return http.StatusNotFound, errors.New("fuel not found")
	}

	data, err := serializeFuel(fuel)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cache.client.HSet(key, strconv.Itoa(fuel.ID), data).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cache *fuelCache) DeleteFuel(id int) (int, error) {
	deleted, err := cache.client.HDel(key, strconv.Itoa(id)).Result()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if deleted == 0 {
		return http.StatusNotFound, errors.New("fuel not found")
	}

	return http.StatusOK, nil
}
