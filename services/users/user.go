package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/vazzuniversebackend/pkg/server/middlewares"
	"github.com/wafi04/vazzuniversebackend/services/auth/sessions"
)

type Users struct {
	MainDB    *sqlx.DB
	ReplicaDB *sqlx.DB
}

func UsersSetUp(db *Users, router *gin.Engine) {
	userRepo := NewUserRepositories(db.MainDB, db.ReplicaDB)
	sessionRepo := sessions.NewSessionRepo(db.MainDB, db.ReplicaDB)
	userService := NewUserServices(userRepo, sessionRepo)
	sessionService := sessions.NewSessionsService(sessionRepo)
	userController := NewUserController(userService, sessionService)
	sessionController := sessions.NewSessionController(sessionService)

	RegisterUserRoutes(router, userController, sessionController)
}

func RegisterUserRoutes(router *gin.Engine, controller *UserController, sessionController *sessions.SessionsController) {
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/register", controller.CreateUser)
		userGroup.POST("/login", controller.Login)
	}
	authGroup := router.Group("/api/users")
	authGroup.Use(middlewares.AuthMiddleware(controller.SessionService))
	{
		// authentication
		authGroup.GET("/profile", controller.GetProfile)
		authGroup.POST("/logout", controller.Logout)
		authGroup.DELETE("/sessions/revoke/:id", sessionController.DeleteSessions)
		authGroup.DELETE("/sessions/revokes", sessionController.ClearSessions)
	}

}
