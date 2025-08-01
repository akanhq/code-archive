package services

import (
	"code_test/message_queue/mongodb_database/blog-system/models"
	"code_test/message_queue/mongodb_database/blog-system/repositories"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthService struct {
	UserRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) Register(username, password, email string) error {
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	return s.UserRepo.Create(&user)
}

func (s *AuthService) BatchCreate(users []*models.User) error {
	return s.UserRepo.BatchCreate(users)

}

func (s *AuthService) Login(username, password string) (string, error) {

	var err error
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil || user.Password != hashPassword(password) {
		return "", fmt.Errorf("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte("secret"))

}
func hashPassword(password string) string {
	return password
}
