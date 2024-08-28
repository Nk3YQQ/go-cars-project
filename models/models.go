package models

import (
	"backend/types"

	"gorm.io/gorm"
)

type RegisterInput struct {
	// Модель регистрации

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type LoginInput struct {
	// Модель входа

	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	// Модель для пользователя

	gorm.Model

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"password"`
}

type Car struct {
	// Модель машин

	gorm.Model

	Brand        string             `json:"brand"`
	Type         string             `json:"type"`
	FuelType     types.Fuel         `json:"fuel_type"`
	Transmission types.Transmission `json:"transmission"`
	Amount       int                `json:"amount"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
}
