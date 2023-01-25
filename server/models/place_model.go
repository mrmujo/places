package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Place struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Longitude 	string             `json:"longitude,omitempty" validate:"required"`
	Latitude  	string             `json:"latitude,omitempty" validate:"required"`
}
