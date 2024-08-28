package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"backend/models"
	"backend/schemas"

	"gorm.io/gorm"
)

func HasCar(w http.ResponseWriter, app *App, id uint) bool {
	var car models.Car

	if err := app.DB.First(&car, id).Error; err == nil {
		return false
	}

	return true
}

func (app *App) CreateCar(w http.ResponseWriter, r *http.Request) {
	// Контроллер для создания машины

	var input models.Car

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := reflect.ValueOf(input)
	typeOfValue := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.String && field.String() == "" {
			http.Error(w, fmt.Sprintf("Missing required field: %s", typeOfValue.Field(i).Name), http.StatusBadRequest)
			return
		}

		if field.Kind() == reflect.Int && field.Int() == 0 {
			http.Error(w, fmt.Sprintf("Missing required field: %s", typeOfValue.Field(i).Name), http.StatusBadRequest)
			return
		}
	}

	userIDValue := r.Context().Value("userID")

	userID := userIDValue.(uint)

	car := models.Car{
		Brand:        input.Brand,
		Type:         input.Type,
		FuelType:     input.FuelType,
		Transmission: input.Transmission,
		Amount:       input.Amount,
		UserID:       userID,
	}

	app.DB.Create(&car)

	carSerializer := schemas.Car{
		Brand:        car.Brand,
		Type:         car.Type,
		FuelType:     car.FuelType,
		Transmission: car.Transmission,
		Amount:       car.Amount,
		UserID:       car.UserID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(carSerializer)
}

func (app *App) GetAllCars(w http.ResponseWriter, r *http.Request) {
	// Чтение всех машин

	var cars []models.Car
	var totalCars int64

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if err := app.DB.Model(&models.Car{}).Count(&totalCars).Error; err != nil {
		http.Error(w, "Could't count cars", http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * pageSize

	if err := app.DB.Limit(pageSize).Offset(offset).Find(&cars).Error; err != nil {
		http.Error(w, "Counld't find retrieve cars", http.StatusInternalServerError)
		return
	}

	serializedCars := make([]schemas.Car, len(cars))

	for i, car := range cars {
		serializedCars[i] = schemas.Car{
			Brand:        car.Brand,
			Type:         car.Type,
			FuelType:     car.FuelType,
			Transmission: car.Transmission,
			Amount:       car.Amount,
			UserID:       car.UserID,
		}
	}

	response := map[string]interface{}{
		"cars":       serializedCars,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int((totalCars + int64(pageSize) - 1) / int64(pageSize)),
	}

	json.NewEncoder(w).Encode(response)
}

func (app *App) GetCarByID(w http.ResponseWriter, r *http.Request) {
	// Получение машины по его ID

	idStr := r.URL.Path[len("/cars/"):]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID for car", http.StatusBadRequest)
		return
	}

	var car models.Car

	if err := app.DB.First(&car, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "car not found", http.StatusNotFound)
		} else {
			http.Error(w, "could't retrieve car", http.StatusInternalServerError)
		}

		return
	}

	serializedCar := schemas.Car{
		Brand:        car.Brand,
		Type:         car.Type,
		FuelType:     car.FuelType,
		Transmission: car.Transmission,
		Amount:       car.Amount,
		UserID:       car.UserID,
	}

	json.NewEncoder(w).Encode(serializedCar)
}

func (app *App) CarsCreateListView(w http.ResponseWriter, r *http.Request) {
	// Котроллер для обработки POST и GET запросов на марштруте /cars

	switch r.Method {
	case http.MethodGet:
		app.GetAllCars(w, r)

	case http.MethodPost:
		app.CreateCar(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
