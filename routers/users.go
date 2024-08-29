package routers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/auth"
	"backend/models"
	"backend/schemas"
	"backend/validators"
)

type App struct {
	DB *gorm.DB
}

func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Контроллер для создания пользователя

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.RegisterInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors := validators.ValidateRegisterForm(input.FirstName, input.LastName, input.Email, input.Password, input.PasswordConfirm)

	for k := range validationErrors {
		if validationErrors[k] != "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validators.ErrorResponse{Details: validationErrors})
			return
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	app.DB.Create(&user)

	userSerializer := schemas.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userSerializer)
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Контроллер для входа в аккаунт

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := validators.ValidateUserLogin(app.DB, input.Email, input.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(user.ID)

	if err != nil {
		http.Error(w, "could not generate token", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"access": token})
}

func (app *App) UserProfile(w http.ResponseWriter, r *http.Request) {
	// Контроллер для профиля текущего пользователя

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDValue := r.Context().Value("userID")

	if userIDValue == nil {
		http.Error(w, "user ID is missing in context", http.StatusUnauthorized)
	}

	userID, ok := userIDValue.(uint)

	if !ok {
		http.Error(w, "user ID has wrong type", http.StatusUnauthorized)
	}

	var user models.User
	app.DB.First(&user, userID)

	userSerializer := schemas.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	json.NewEncoder(w).Encode(userSerializer)
}
