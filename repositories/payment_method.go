package repositories

import (
	"context"
	"net/http"
	"time"

	"carbon-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentMethodRepository interface {
	GetAllPaymentMethods() ([]models.PaymentMethod, int, error)
	GetPaymentMethodByCode(code string) (models.PaymentMethod, int, error)
}

type paymentMethodRepository struct {
	MongoCollection *mongo.Collection
}

func NewPaymentMethodRepository(mc *mongo.Collection) PaymentMethodRepository {
	return &paymentMethodRepository{mc}
}

func (pmr *paymentMethodRepository) GetAllPaymentMethods() ([]models.PaymentMethod, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := pmr.MongoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer cursor.Close(ctx)

	var paymentMethods []models.PaymentMethod
	for cursor.Next(ctx) {
		var paymentMethod models.PaymentMethod
		if err := cursor.Decode(&paymentMethod); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		paymentMethods = append(paymentMethods, paymentMethod)
	}

	return paymentMethods, http.StatusOK, nil
}

func (pmr *paymentMethodRepository) GetPaymentMethodByCode(code string) (models.PaymentMethod, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var paymentMethod models.PaymentMethod
	err := pmr.MongoCollection.FindOne(ctx, bson.M{"code": code}).Decode(&paymentMethod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.PaymentMethod{}, http.StatusNotFound, err
		}
		return models.PaymentMethod{}, http.StatusInternalServerError, err
	}

	return paymentMethod, http.StatusOK, nil
}
