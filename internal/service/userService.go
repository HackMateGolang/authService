package service

import (
	"context"
	"time"

	"github.com/DestWish/HackMate/Auth-service/internal/models"
	"github.com/DestWish/HackMate/Auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secret string = "HackMate_secretKey"

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser (ctx context.Context, req *models.UserCreateRequest) (string, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := &models.User{Login: req.Login, Email: req.Email, PasswordHash: passwordHash}

	login, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := generateToken(login)
	if err != nil {
		return "", err 
	}

	return token, nil
}

func generateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

