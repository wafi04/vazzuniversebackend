package users

import (
	"context"
	"strings"

	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
)

type UserService struct {
	userRepo *UserRepositories
}

func NewUserServices(userepo *UserRepositories) *UserService {
	return &UserService{userRepo: userepo}
}

type UserServices interface {
	Create(ctx context.Context, req *CreateUser) (*UserData, *response.ResponseError)
	GetUserByID(ctx context.Context, userID string) (*UserData, error)
}

func (us *UserService) Create(ctx context.Context, req *CreateUser) (*UserData, *response.ResponseError) {
	// Validate required fields
	if req.Username == "" {
		return nil, response.NewResponseError(400, ErrMissingField, "Username is required")
	}

	if req.Email == "" {
		return nil, response.NewResponseError(400, ErrMissingField, "Email is required")
	}

	if req.Password == nil || *req.Password == "" {
		return nil, response.NewResponseError(400, ErrMissingField, "Password is required")
	}

	// Check if username already exists
	_, err := us.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if !isNotFoundError(err) {
			return nil, ErrInternalServerError(err)
		}
	} else {
		return nil, ErrUsernameTakenError(req.Username)
	}

	_, err = us.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !isNotFoundError(err) {
			return nil, ErrInternalServerError(err)
		}
	} else {
		return nil, ErrEmailTakenError(req.Email)
	}

	return us.userRepo.Create(ctx, *req)
}

func (us *UserService) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
	return us.userRepo.GetUserByID(ctx, userID)
}

func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
