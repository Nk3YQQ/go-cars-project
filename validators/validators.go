package validators

import (
	"errors"
	"regexp"

	"backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ValidatePassword(password, PasswordConfirm string) error {
	// Валидация пароля при регистрации

	if len(password) < 6 {
		return errors.New("длина пароля не может быть меньше шести символов")
	}

	if password != PasswordConfirm {
		return errors.New("пароли не совпаают")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString

	if !hasDigit(password) {
		return errors.New("пароль должен состоять из букв и цифр")
	}

	return nil
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
