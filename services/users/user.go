package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
	userController := NewUserController(userService)

	RegisterUserRoutes(router, userController)
}

func RegisterUserRoutes(router *gin.Engine, controller *UserController) {
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/register", controller.CreateUser)
		userGroup.POST("/login", controller.Login)

	}
}
