package repository

import (
	"context"
	"fmt"
	// "gorm.io/datatypes"
	// "time"

	"github.com/HackMateGolang/AuthService/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	db	*gorm.DB
	redisClient *redis.Client
}


func newUserRepo(db *gorm.DB, redisClient *redis.Client) *UserRepository {
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


func (r* UserRepository) GetUserLogin(ctx context.Context, user *models.User) (models.User, error) {
	key := userCacheKey(user.Login)
	var newUser models.User

	err := r.redisClient.HGetAll(ctx, key).Scan(&newUser)
	if err == nil && newUser.Login != "" {
		return newUser, nil
	}
	
	if err := r.db.Where("Login= ?", user.Login).First(&newUser).Error; err != nil {
		return newUser, fmt.Errorf("Repository: User not found: %w", err)
	}
	return newUser, r.userCaching(ctx, &newUser)
} 

func (r* UserRepository) GetUserEmail(ctx context.Context, user *models.User) (models.User, error) {
	key := userCacheKey(user.Email)
	var newUser models.User

	err := r.redisClient.HGetAll(ctx, key).Scan(&newUser)
	if err == nil && newUser.Email != "" {
		return newUser, nil
	}
	
	if err := r.db.Where("Email= ?", user.Email).First(&newUser).Error; err != nil {
		return newUser, fmt.Errorf("Repository: User not found: %w", err)
	}
	return newUser, r.userCaching(ctx, &newUser)
} 

func (r* UserRepository) PatchUser(ctx context.Context, req *models.UserPatchRequest) (bool, error){
	modifiedUserLogin := r.db.Where("Login = ?", req.Login).Model(&models.User{}).Updates(req)
	modifiedUserEmail := r.db.Where("Email = ?", req.Email).Model(&models.User{}).Updates(req)

	if modifiedUserLogin.Error != nil {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserLogin.Error)
	}

	if modifiedUserEmail.Error != nil {
		return  false, fmt.Errorf("Repository: User not found: %w", modifiedUserEmail.Error)
	}
	
	if modifiedUserLogin.RowsAffected == 0 {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserLogin.Error)
	}
	if modifiedUserEmail.RowsAffected == 0 {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserEmail.Error)
	}

	var patchedUser models.User
	if err := r.db.Where("Login=?", req.Login).First(&patchedUser).Error; err != nil {
		return false, fmt.Errorf("Repository: Patched user not found: %w", err)
	}
	if err := r.db.Where("Email=?", req.Email).First(&patchedUser).Error; err != nil {
		return false, fmt.Errorf("Repository: Patched user not found: %w", err)
	}

	return true, r.userCaching(ctx, &patchedUser)

}
func (r* UserRepository) PutUser(ctx context.Context, req *models.UserPutRequest) (bool, error){
	modifiedUserLogin := r.db.Model(&models.User{}).Where("Login = ?", req.Login).Select("*").Updates(req)
	modifiedUserEmail := r.db.Model(&models.User{}).Where("Email = ?", req.Email).Select("*").Updates(req)

	
	if modifiedUserLogin.Error != nil {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserLogin.Error)
	}

	if modifiedUserEmail.Error != nil {
		return  false, fmt.Errorf("Repository: User not found: %w", modifiedUserEmail.Error)
	}

	if modifiedUserLogin.RowsAffected == 0 {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserLogin.Error)
	}

	if modifiedUserEmail.RowsAffected == 0 {
		return false, fmt.Errorf("Repository: User not found: %w", modifiedUserEmail.Error)
	}

	var putUser models.User
	if err := r.db.Where("Login=?", req.Login).First(&putUser).Error; err != nil {
		return false, fmt.Errorf("Repository: Patched user not found: %w", err)
	}
	if err := r.db.Where("Email=?", req.Email).First(&putUser).Error; err != nil {
		return false, fmt.Errorf("Repository: Patched user not found: %w", err)
	}
	return true, r.userCaching(ctx, &putUser)

}



func userCacheKey(userLogin string) string {
	return fmt.Sprintf("user:%v", userLogin)
}