package repositories

import (
	"context"
	"net/http"
	"time"

	"carbon-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentMethodRepository interface {
	GetAllPaymentMethods() ([]models.PaymentMethod, int, error)
	GetPaymentMethodByCode(code string) (models.PaymentMethod, int, error)
	CreatePaymentMethod(paymentMethod models.PaymentMethod) (models.PaymentMethod, int, error)
	UpdatePaymentMethod(id string, paymentMethod models.PaymentMethod) (models.PaymentMethod, int, error)
	DeletePaymentMethod(id string) (int, error)
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

	if len(paymentMethods) == 0 {
		paymentMethods = []models.PaymentMethod{}
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

func (pmr *paymentMethodRepository) CreatePaymentMethod(paymentMethod models.PaymentMethod) (models.PaymentMethod, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := pmr.MongoCollection.InsertOne(ctx, paymentMethod)
	if err != nil {
		return models.PaymentMethod{}, http.StatusInternalServerError, err
	}

	paymentMethod.ID = result.InsertedID.(primitive.ObjectID)

	return paymentMethod, http.StatusCreated, nil
}

func (pmr *paymentMethodRepository) UpdatePaymentMethod(id string, paymentMethod models.PaymentMethod) (models.PaymentMethod, int, error) {
	paymentMethodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.PaymentMethod{}, http.StatusBadRequest, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = pmr.MongoCollection.ReplaceOne(ctx, bson.M{"_id": paymentMethodID}, paymentMethod)
	if err != nil {
		return models.PaymentMethod{}, http.StatusInternalServerError, err
	}

	paymentMethod.ID = paymentMethodID

	return paymentMethod, http.StatusOK, nil
}

func (pmr *paymentMethodRepository) DeletePaymentMethod(id string) (int, error) {
	paymentMethodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = pmr.MongoCollection.DeleteOne(ctx, bson.M{"_id": paymentMethodID})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
