package users

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/wafi04/vazzuniversebackend/pkg/server/middlewares"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/generate"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"github.com/wafi04/vazzuniversebackend/services/auth/sessions"
)

type UserService struct {
	userRepo    *UserRepositories
	sessionRepo *sessions.SessionRepo
}

func NewUserServices(userepo *UserRepositories, sessionRepo *sessions.SessionRepo) *UserService {
	return &UserService{userRepo: userepo, sessionRepo: sessionRepo}
}

type UserServices interface {
	Create(ctx context.Context, req *CreateUser) (*UserData, *response.ResponseError)
	GetUserByID(ctx context.Context, userID string) (*UserData, error)
	LoginWithSession(ctx context.Context, req *LoginUser, ipAddress, userAgent, deviceInfo string) (*UserData, *sessions.SessionsData, error)
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

func (us *UserService) LoginWithSession(ctx context.Context, req *LoginUser, ipAddress, userAgent, deviceInfo string) (*UserData, *sessions.SessionsData, *response.ResponseError) {
	// First authenticate the user
	userData, err := us.userRepo.Login(ctx, req)
	if err != nil {
		switch err.Error() {
		case string(ErrInvalidCredentials):
			return nil, nil, response.NewResponseError(
				http.StatusUnauthorized,
				"Invalid credentials",
				"The username you entered doesn't exist",
			)
		case string(ErrInvalidPassword):
			return nil, nil, response.NewResponseError(
				http.StatusUnauthorized,
				"Invalid password",
				"The password you entered is incorrect",
			)
		default:
			return nil, nil, response.NewResponseError(
				http.StatusInternalServerError,
				"Internal server error",
				"An unexpected error occurred during login",
			)
		}
	}
	sessionID := generate.GenerateRandomID(&generate.IDOpts{
		Amount: 10,
	})

	accessToken, err := middlewares.GenerateToken(&middlewares.UserData{
		UserID:    userData.UserID,
		Username:  userData.Username,
		Fullname:  userData.FullName,
		Email:     userData.Email,
		Role:      string(userData.Role),
		Balance:   userData.Balance,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		SessionID: sessionID,
	}, 3600*24)
	if err != nil {
		return nil, nil, response.NewResponseError(
			http.StatusBadGateway,
			"Invalid Generate Token",
			"BAD_GATEWAY",
		)
	}

	sessionReq := &sessions.CreateSession{
		SessionID:    sessionID,
		UserID:       userData.UserID,
		AccessToken:  accessToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		DeviceInfo:   deviceInfo,
		LastActivity: time.Now(),
		IsAccess:     true,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	session, err := us.sessionRepo.Create(ctx, sessionReq)
	if err != nil {
		return nil, nil, response.NewResponseError(
			http.StatusBadRequest,
			"Faield To Create Sessions",
			sessions.ErrSessionsInvalid,
		)
	}

	return userData, session, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
	return us.userRepo.GetUserByID(ctx, userID)
}

func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (us *UserService) Logout(ctx context.Context, userID string) error {
	return us.userRepo.Logout(ctx, userID)
}
