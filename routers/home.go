package routers

import (
	"encoding/json"
	"net/http"

	"backend/models"
)

type HomePage struct {
	Greetings      string       `json:"greetings"`
	TotalCarCount  int64        `json:"car_count"`
	TotalUserCount int64        `json:"user_count"`
	RecomendedCars []models.Car `json:"recomended_cars"`
}

func (app *App) HomePageView(w http.ResponseWriter, r *http.Request) {
	// Контроллер для главной страницы

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var totalCarsCount int64
	var totalUserCount int64
	var recomendedCars []models.Car

	greeting := "Добро пожаловать на сраницу магазина машин E&I Cars! Здесь Вы сможете найти ту самую машину, которая подойдёт именно Вам!"

	app.DB.Model(&models.Car{}).Count(&totalCarsCount)
	app.DB.Model(&models.User{}).Count(&totalUserCount)
	app.DB.Limit(10).Find(&recomendedCars)

	homePage := HomePage{
		Greetings:      greeting,
		TotalCarCount:  totalCarsCount,
		TotalUserCount: totalUserCount,
		RecomendedCars: recomendedCars,
	}

	json.NewEncoder(w).Encode(&homePage)
}
