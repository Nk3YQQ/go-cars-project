package schemas

import (
	"backend/types"
)

type Car struct {
	// Сериализатор для машины

	Brand        string             `json:"brand"`
	Type         string             `json:"type"`
	FuelType     types.Fuel         `json:"fuel_type"`
	Transmission types.Transmission `json:"transmission"`
	Amount       int                `json:"amount"`
	UserID       uint               `json:"user_id"`
}

type User struct {
	// Сериализатор для пользователя

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" json:"email"`
}
