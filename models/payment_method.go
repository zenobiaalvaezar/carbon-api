package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethod struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code string             `bson:"code" json:"code"`
}
