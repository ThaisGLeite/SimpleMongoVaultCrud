package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty" validate:"omitempty,min=2"`     // Ensure a minimum length of 2
	Age      int                `bson:"age,omitempty" validate:"omitempty,gte=1"`      // Ensure age is non-negative and at leat 1
	Email    string             `bson:"email,omitempty" validate:"omitempty,email"`    // Validate if it is a proper email format
	Password string             `bson:"password,omitempty" validate:"omitempty,min=8"` // Ensure a minimum length of 8
	Address  string             `bson:"address,omitempty" validate:"omitempty,min=5"`  // Ensure a minimum length of 5
}
