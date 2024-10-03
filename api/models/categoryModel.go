package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Name  string               `json:"name" bson:"name"`
	Todos []primitive.ObjectID `json:"todos" bson:"todos"`
}
