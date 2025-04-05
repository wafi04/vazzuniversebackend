package users

import (
	"context"
	"errors"
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

func (us *UserService) Create(ctx context.Context, req *CreateUSer) (*UserData, error) {
    // Validasi input jika diperlukan
    if req.Username == "" || req.Email == ""  {
        return nil, errors.New("username, email, and password are required")
    }

    // Panggil repository untuk menyimpan data user
    return us.userRepo.Create(ctx, *req)
}

// GetUserByID retrieves a user by their ID
func (us *UserService) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
    // Panggil repository untuk mendapatkan data user
    return us.userRepo.GetUserByID(ctx, userID)
}