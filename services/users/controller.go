package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/vazzuniversebackend/pkg/server/middlewares"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"github.com/wafi04/vazzuniversebackend/services/auth/sessions"
)

type UserController struct {
	UserService *UserService
}

func NewUserController(UserService *UserService) *UserController {
	return &UserController{UserService: UserService}
}

type ReqData struct {
	FullName *string `json:"fullName"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
	WhatsApp string  `json:"whatsapp"`
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

	user, err := uc.UserService.userRepo.GetUserByEmail(ctx, users.UserID)

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request",
			"error":   err.Error(),
		})
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

	user, session, err := uc.UserService.userRepo.LoginWithSession(
		ctx.Request.Context(),
		&loginReq,
		(*sessions.SessionRepo)(uc.UserService.userRepo),
		clientIP,
		userAgent,
		deviceInfo,
	)

	if err != nil {
		respErr := response.NewResponseError(
			http.StatusUnauthorized,
			ErrUnauthorized,
			"User Not Authorized",
		)
		response.Error(ctx, respErr)
		return
	}

	// Set cookie with access token (optional)
	middlewares.SetTokenCookie(ctx, "auth_token", session.AccessToken, session.ExpiresAt.Minute())
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
