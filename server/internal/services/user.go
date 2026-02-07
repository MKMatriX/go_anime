package services

import (
	"errors"
	"fmt"
	"go_anime/internal/models"
	"go_anime/internal/requests"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) RegisterUser(request *requests.UserRegisterRequest) (*models.UserModel, error) {
	// hash the password
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("User registration error")
	}

	createdUser := models.UserModel{
		Login:    request.Login,
		Password: hashedPassword,
	}
	result := s.db.Create(&createdUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("User registration error")
	}

	return &createdUser, nil
}

func (s *UserService) LoginUser(request *requests.UserLoginRequest) (*models.UserModel, error) {
	user, err := s.GetUserByLogin(request.Login)
	if err != nil {
		return nil, errors.New("Login and/or password are incorrect")
	}

	if !checkPasswordHash(request.Password, user.Password) {
		fmt.Println("password mismatch")
		return nil, errors.New("Login and/or password are incorrect")
	}

	return user, nil
}

func (s *UserService) GetUserByLogin(login string) (*models.UserModel, error) {
	var user models.UserModel
	result := s.db.Where("login = ? ", login).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
