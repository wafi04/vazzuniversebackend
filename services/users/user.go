package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/vazzuniversebackend/pkg/server/middlewares"
)

type Users struct {
	MainDB    *sqlx.DB
	ReplicaDB *sqlx.DB
}

func UsersSetUp(db *Users, router *gin.Engine) { // Add router parameter
	userRepo := NewUserRepositories(db.MainDB, db.ReplicaDB)
	userService := NewUserServices(userRepo)
	userController := NewUserController(userService)

	// Register routes
	RegisterUserRoutes(router, userController)
}

func RegisterUserRoutes(router *gin.Engine, controller *UserController) {
	userGroup := router.Group("/api/users")
	authGroup := userGroup.Use(middlewares.AuthMiddleware())
	{
		// userGroup.GET("/", controller.GetAllUsers)
		authGroup.GET("/:id", controller.GetProfile)
		userGroup.POST("/", controller.CreateUser)
		// userGroup.PUT("/:id", controller.UpdateUser)
		// userGroup.DELETE("/:id", controller.DeleteUser)

		// // Add more routes as needed
		// userGroup.POST("/login", controller.Login)
		// userGroup.POST("/logout", controller.Logout)
	}
}
