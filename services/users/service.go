package users

import (
	"context"
	"errors"

	"github.com/wafi04/vazzuniversebackend/pkg/utils"
)

type UserService struct {
	userRepo  *UserRepositories
}

func NewUserServices(userepo *UserRepositories) *UserService{
	return &UserService{userRepo: userepo}
}


type  UserServices interface {
	create (ctx context.Context,req *CreateUSer) (*UserData,error)
}

func (us *UserService) Create(ctx context.Context, req *CreateUSer) (*UserData, *utils.ResponseError) {
 	if req.Username == "" {
		return nil, ErrUsernameTakenError(req.Username)
	}
	if req.Email == "" {
		return nil, ErrEmailTakenError(req.Email)
	}
	
	existingUser, err := us.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err,errors.New("not found")) {
		} else {
			return nil, ErrInternalServerError(err)
		}
	} else if existingUser != nil {
		// User exists
		return nil, ErrUsernameTakenError(req.Username)
	}

	existingEmail, err := us.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
		} else {
			return nil, ErrInternalServerError(err)
		}
	} else if existingEmail != nil {
		return nil, ErrEmailTakenError(req.Email)
	}

    return us.userRepo.Create(ctx, *req)
}

func (us *UserService) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
    return us.userRepo.GetUserByID(ctx, userID)
}