package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/vazzuniversebackend/pkg/server/middlewares"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"github.com/wafi04/vazzuniversebackend/services/auth/sessions"
)

type UserController struct {
	UserService    *UserService
	SessionService *sessions.SessionService
}

func NewUserController(UserService *UserService, sessionService *sessions.SessionService) *UserController {
	return &UserController{UserService: UserService, SessionService: sessionService}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var req ReqData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respErr := response.NewResponseError(http.StatusBadRequest, ErrInvalidInput, "Invalid Format request").WithDetails(err.Error())
		response.Error(ctx, respErr)
		return
	}

	if req.Username == "" {
		respErr := response.NewResponseError(
			http.StatusBadRequest,
			ErrMissingField,
			"Username is requeired",
		)
		response.Error(ctx, respErr)
		return
	}
	if req.Email == "" {
		respErr := response.NewResponseError(
			http.StatusBadRequest,
			ErrMissingField,
			"Email is requeired",
		)
		response.Error(ctx, respErr)
		return
	}

	if req.Password == nil || *req.Password == "" {
		respErr := response.NewResponseError(
			http.StatusBadRequest,
			ErrMissingField,
			"Password is required",
		)
		response.Error(ctx, respErr)
		return
	}

	userData, respErr := uc.UserService.Create(ctx, &CreateUser{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "Member",
		Balance:  0,
	})
	if respErr != nil {
		log.Printf("%s", respErr.Message)
		response.Error(ctx, respErr)
		return
	}

	response.Success(ctx, http.StatusCreated, userData)

}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	users, err := middlewares.GetUserFromGinContext(ctx)

	if err != nil {
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUnauthorized,
			"Unauthorized",
		)
		response.Error(ctx, respErr)
		return
	}

	user, err := uc.UserService.userRepo.GetUserByUsername(ctx, users.Username)

	if err != nil {
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUserNotFound,
			"User Not Found,Please Login Frist!",
		)
		response.Error(ctx, respErr)
		return
	}

	response.Success(ctx, http.StatusOK, user)
}

func (uc *UserController) Login(ctx *gin.Context) {
	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	deviceInfo := fmt.Sprintf("Device accessing from %s", clientIP)

	var loginReq LoginUser
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		respErr := response.NewResponseError(
			http.StatusBadRequest,
			ErrInvalidInput,
			"Invalid Format Request",
		)
		response.Error(ctx, respErr)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		respErr := response.NewResponseError(
			http.StatusBadRequest,
			ErrInvalidCredentials,
			"Invalid Login Username And Password Is Required",
		)
		response.Error(ctx, respErr)
		return
	}

	user, session, err := uc.UserService.LoginWithSession(
		ctx.Request.Context(),
		&loginReq,
		clientIP,
		userAgent,
		deviceInfo,
	)

	if err != nil {
		log.Printf("Login failed: %v", err)
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUnauthorized,
			"User Not Authorized",
		)
		response.Error(ctx, respErr)
		return
	}
	log.Printf("second : %d", session.ExpiresAt.Second())
	log.Printf("day : %d", session.ExpiresAt.Day())
	log.Printf("hor : %d", session.ExpiresAt.Hour())

	middlewares.SetTokenCookie(ctx, "auth_token", session.AccessToken, 24*60*60)
	ctx.Header("Authorization", "Bearer "+session.AccessToken)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"user":       user,
			"token":      session.AccessToken,
			"session_id": session.SessionID,
			"expires_at": session.ExpiresAt,
		},
	})

}

func (uc *UserController) Logout(ctx *gin.Context) {
	userData, err := middlewares.GetUserFromGinContext(ctx)
	if err != nil {
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUnauthorized,
			"Unauthorized",
		)
		response.Error(ctx, respErr)
		return
	}

	uc.UserService.userRepo.Logout(ctx, userData.UserID)
	uc.UserService.userRepo.DeleteSession(ctx, userData.SessionID)

	middlewares.ClearTokens(ctx)
	response.Success(ctx, http.StatusOK, "Logout successful")
}

func (uc *UserController) ClearSessions(ctx *gin.Context) {
	userData, err := middlewares.GetUserFromGinContext(ctx)
	if err != nil {
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUnauthorized,
			"Unauthorized",
		)
		response.Error(ctx, respErr)
		return
	}

	err = uc.SessionService.InvalidateAllUserSessions(ctx, userData.UserID)
	if err != nil {
		respErr := response.NewResponseError(
			http.StatusInternalServerError,
			"Failed to clear user sessions",
			"INTERNAL_SERVER_ERROR",
		)
		response.Error(ctx, respErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "All user sessions cleared successfully",
	})
}
