package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system. The struct includes fields for the user's
// ID, name, age, email, password, and address. The bson tags tell the mongo-driver how to name
// the properties in BSON when it marshals the data to be stored in MongoDB. The validate tags
// are used to set constraints on the data when it's being validated.
type User struct {
	// ID is the unique identifier of the user, represented as a MongoDB ObjectID
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// Name of the user, must be at least 2 characters long
	Name string `bson:"name,omitempty" validate:"omitempty,min=2"`

	// Age of the user, must be at least 1
	Age int `bson:"age,omitempty" validate:"omitempty,gte=1"`

	// Email of the user, must be a valid email format
	Email string `bson:"email,omitempty" validate:"omitempty,email"`

	// Password of the user, must be at least 8 characters long
	Password string `bson:"password,omitempty" validate:"omitempty,min=8"`

	// Address of the user, must be at least 5 characters long
	Address string `bson:"address,omitempty" validate:"omitempty,min=5"`
}
