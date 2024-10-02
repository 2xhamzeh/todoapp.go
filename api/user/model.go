package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
}

type AuthDTO struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
