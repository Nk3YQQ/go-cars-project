package validators

import (
	"errors"
	"regexp"

	"backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	// Структура для ответа сервера на обработку ошибок

	Details map[string]string `json:"details"`
}

func ValidateRegisterForm(firstName, lastName, email, password, PasswordConfirm string) map[string]string {
	// Валидация пароля при регистрации

	validationErrors := make(map[string]string)

	if firstName == "" {
		validationErrors["first_name"] = "Поле 'first_name' обязательно"
	}

	if lastName == "" {
		validationErrors["last_name"] = "Поле 'last_name' обязательно"
	}

	if email == "" {
		validationErrors["email"] = "Поле 'email' обязательно"
	} else {
		if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
			validationErrors["email"] = "Формат поля 'email' некорректный"
		}
	}

	if len(password) < 6 {
		validationErrors["password"] = "Длина пароля не может быть меньше шести символов"
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString

	if !hasDigit(password) {
		validationErrors["password"] = "Пароль должен состоять из букв и цифр"
	}

	if password != PasswordConfirm {
		validationErrors["passwordConfirm"] = "Пароли не совпадают"
	}

	return validationErrors
}

func ValidateUserLogin(db *gorm.DB, email, password string) (*models.User, error) {
	// Валидация входа пользователя

	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("неверный логин или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("неверный логин или пароль")
	}

	return &user, nil
}
