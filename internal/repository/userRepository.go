package repository

import (
	"context"
	"fmt"
	// "gorm.io/datatypes"
	// "time"

	"github.com/DestWish/HackMate/Auth-service/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	db	*gorm.DB
	redisClient *redis.Client
}


func NewUserRepo(db *gorm.DB, redisClient *redis.Client) *UserRepository {
	return &UserRepository{db: db, redisClient:  redisClient}
}


func (r *UserRepository) userCaching(ctx context.Context, user *models.User) error {
	key := userCacheKey(user.Login)

	if err := r.redisClient.HSet(ctx, key, user).Err(); err != nil {
		return fmt.Errorf("Repo: Cache failed: %w", err)
	}

	return nil
}


func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	user.IsVerified = false
	user.Role = "User"
	
	if err := r.db.Model(&models.User{}).Create(user).Error; err != nil {
		return "...nothing...", fmt.Errorf("Repository: Create user failed! %w", err)
	}
	return user.Login, r.userCaching(ctx, user)
}


func userCacheKey(userLogin string) string {
	return fmt.Sprintf("user:%v", userLogin)
}