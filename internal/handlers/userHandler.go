package handlers

import (
	"context"

	userpb "github.com/HackMateGolang/AuthService/api/proto/v1"
	"github.com/HackMateGolang/AuthService/internal/models"
	"github.com/HackMateGolang/AuthService/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	userpb.UnimplementedAuthServiceServer
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := h.service.CreateUser(ctx, &models.UserCreateRequest{Login: req.Login, Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		Token: user,
	}, nil
}


func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, &models.UserGetRequest{Login: req.Login, Email: req.Email})
	if err != nil {
		return nil, err
	}
	return  &userpb.GetUserResponse{
		Login: user.Login,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified: user.IsVerified,
		Role: user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
	},nil
}
func (h *UserHandler) PatchUser(ctx context.Context, req *userpb.PatchUserRequest) (*userpb.PatchUserResponse, error) {
	createdAtTime := req.CreatedAt.AsTime()
	ok, err := h.service.PatchUser(ctx, &models.UserPatchRequest{
		Login: req.Login,
		Email: req.Email,
		PasswordHash: req.PasswordHash,
		IsVerified: req.IsVerified,
		Role: req.Role,
		CreatedAt: &createdAtTime,
	})
		if err != nil {
		return nil, err
	}
	return &userpb.PatchUserResponse{Ok: ok}, nil
}
func (h *UserHandler) PutUser(ctx context.Context, req *userpb.PutUserRequest) (*userpb.PutUserResponse, error) {
	createdAtTime := req.CreatedAt.AsTime()
	ok, err := h.service.PutUser(ctx, &models.UserPutRequest{
		Login: req.Login,
		Email: req.Email,
		PasswordHash: req.PasswordHash,
		IsVerified: req.IsVerified,
		Role: req.Role,
		CreatedAt: createdAtTime,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.PutUserResponse{Ok: ok}, nil
}
func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	ok, err := h.service.DeleteUser(ctx, &models.UserDeleteRequest{
		Login: req.Login,
		Email: req.Email,
	})
	if err != nil {
		return &userpb.DeleteUserResponse{Ok: ok}, err
	}
	return &userpb.DeleteUserResponse{Ok: ok}, nil
}
