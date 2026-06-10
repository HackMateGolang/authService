package handlers

import (
	"context"

	authpb "github.com/HackMateGolang/proto-contracts/gen/go/auth/v1"
	"github.com/HackMateGolang/AuthService/internal/models"
	"github.com/HackMateGolang/AuthService/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	authpb.UnimplementedAuthServiceServer
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	user, err := h.service.CreateUser(ctx, &models.UserCreateRequest{Login: req.Login, Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, err
	}

	return &authpb.CreateUserResponse{
		Token: user,
	}, nil
}


func (h *UserHandler) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, &models.UserGetRequest{Login: req.Login, Email: req.Email})
	if err != nil {
		return nil, err
	}
	return  &authpb.GetUserResponse{
		Login: user.Login,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified: user.IsVerified,
		Role: user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
	},nil
}
func (h *UserHandler) PatchUser(ctx context.Context, req *authpb.PatchUserRequest) (*authpb.PatchUserResponse, error) {
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
	return &authpb.PatchUserResponse{Ok: ok}, nil
}
func (h *UserHandler) PutUser(ctx context.Context, req *authpb.PutUserRequest) (*authpb.PutUserResponse, error) {
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

	return &authpb.PutUserResponse{Ok: ok}, nil
}
func (h *UserHandler) DeleteUser(ctx context.Context, req *authpb.DeleteUserRequest) (*authpb.DeleteUserResponse, error) {
	ok, err := h.service.DeleteUser(ctx, &models.UserDeleteRequest{
		Login: req.Login,
	})
	if err != nil {
		return &authpb.DeleteUserResponse{Ok: ok}, err
	}
	return &authpb.DeleteUserResponse{Ok: ok}, nil
}
