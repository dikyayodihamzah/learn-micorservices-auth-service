package domain

import "time"

type ResetPasswordToken struct {
	Tokens    string    `json:"tokens" bson:"tokens"`
	Email     string    `json:"email" bson:"email"`
	Phone     string    `json:"phone" bson:"phone"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
