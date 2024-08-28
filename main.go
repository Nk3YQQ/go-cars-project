package main

import (
	"log"
	"net/http"

	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/routers"
)

func main() {
	// Функция для запуска приложения

	db := database.InitDB()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	app := &routers.App{DB: db}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Car{})

	http.HandleFunc("/register", app.RegisterUser)
	http.HandleFunc("/login", app.LoginUser)
	http.HandleFunc("/profile", auth.Authenticate(app.UserProfile))

	http.HandleFunc("/cars", auth.Authenticate(app.CarsCreateListView))
	http.HandleFunc("/cars/{id}", auth.Authenticate(app.GetCarByID))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
