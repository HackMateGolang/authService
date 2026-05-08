package models

import (
	"time"
)

type User struct {
	Login string
	Email string
	PasswordHash string
	IsVerified bool
	Role string
	CreatedAt time.Time
}


type UserCreateRequest struct {
	Login string
	Email string
	Password string
} 

type UserGetRequest struct {
	Login string
	Email string
} 
type UserPatchRequest struct {
	Login string
	Email string
	PasswordHash *string
	IsVerified *bool
	Role *string
	CreatedAt *time.Time
}
type UserPutRequest struct {
	Login string
	Email string
	PasswordHash string
	IsVerified bool
	Role string
	CreatedAt time.Time
}
type UserDeleteRequest struct {
	Login string
	Email string
}