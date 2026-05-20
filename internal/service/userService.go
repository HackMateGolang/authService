package service

import (
	"context"
	"fmt"
	"time"

	"github.com/HackMateGolang/AuthService/internal/models"
	"github.com/HackMateGolang/AuthService/internal/repository"
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

func (s *UserService) CreateUser(ctx context.Context, req *models.UserCreateRequest) (string, error) {
 passwordHash, err := hashPassword(req.Password)
 if err != nil {
  return "", err
 }

 user := &models.User{Login: req.Login, Email: req.Email, PasswordHash: passwordHash}

 if _, err := s.repo.GetUserLogin(ctx, user); err == nil {
  return "Error", fmt.Errorf("Service: User with same login already exists")
 }
 if _, err := s.repo.GetUserEmail(ctx, user); err == nil {
  return "Error", fmt.Errorf("Service: User with same email already exists")
 }

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



func (s *UserService) GetUser(ctx context.Context, req *models.UserGetRequest) (models.User, error) {
	user := &models.User{Login: req.Login, Email: req.Email}
	if _, err := s.repo.GetUserLogin(ctx, user); err != nil{
		return *user, fmt.Errorf("Service: There is no User with this login.") 
	}
	if _, err := s.repo.GetUserEmail(ctx, user); err != nil{
		return *user, fmt.Errorf("Service: There is no User with this email.") 
	}
	FoundUser, err := s.repo.GetUserEmail(ctx, user)
	if err != nil {
		return FoundUser, err
	}
	return FoundUser, nil
}

func (s *UserService) PatchUser(ctx context.Context, req *models.UserPatchRequest) (bool, error) {
	return s.repo.PatchUser(ctx, req)
}
func (s *UserService) PutUser(ctx context.Context, req *models.UserPutRequest) (bool, error) {
	return s.repo.PutUser(ctx, req)
}
func (s *UserService) DeleteUser(ctx context.Context, req *models.UserDeleteRequest) (bool, error) {
	return s.repo.DeleteUser(ctx, req)
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

