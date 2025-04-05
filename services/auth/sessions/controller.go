package sessions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
)

type SessionsController struct {
	sessionService *SessionService
}

func NewSessionController(sessionService *SessionService) *SessionsController {
	return &SessionsController{
		sessionService: sessionService,
	}
}

func (sc *SessionsController) DeleteSessions(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Session ID is required",
		})
		return
	}

	err := sc.sessionService.InvalidateSession(ctx, sessionID)
	if err != nil {
		respErr := response.NewResponseError(
			http.StatusNotFound,
			response.ErrorCode(ErrSessionsInvalid),
			"Session Invalid",
		)
		response.Error(ctx, respErr)
		return
	}
	Success := true
	response.Success(ctx, 200, Success)
}

func (sc *SessionsController) ClearSessions(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ID is required",
		})
		return
	}

	err := sc.sessionService.InvalidateAllUserSessions(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to clear user sessions",
			"error":   err.Error(),
		})
		respErr := response.NewResponseError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
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
