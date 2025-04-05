package sessions

import "github.com/gin-gonic/gin"

type SessionsController struct {
	sessionService *SessionService
}

func NewSessionController(sessionService *SessionService) *SessionsController {
	return &SessionsController{
		sessionService: sessionService,
	}
}

func (sc *SessionsController) CreateSessions(ctx *gin.Context) {

}
