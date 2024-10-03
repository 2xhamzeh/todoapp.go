package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDo struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `bson:"title" json:"title"`
	Text  string             `bson:"text" json:"text"`
	Done  bool               `bson:"done" json:"done"`
}

type CreateTodoDTO struct {
	Title string `json:"title" bson:"title"`
	Text  string `json:"text" bson:"text"`
}

// using this we can update a todo
// either update the text and title fields or change it to done/undone
type UpdateTodoDTO struct {
	Title *string `json:"title,omitempty" bson:"title,omitempty"`
	Text  *string `json:"text,omitempty" bson:"text,omitempty"`
	Done  *bool   `json:"done,omitempty" bson:"done,omitempty"`
}
