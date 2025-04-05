package server

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wafi04/vazzuniversebackend/pkg/config"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"github.com/wafi04/vazzuniversebackend/services/users"
)

func SetUpAllRoutes() {
	log := response.NewLogger()

	if err := config.LoadConfig("local"); err != nil {
		log.Log(response.ErrorLevel, "Error loading config: %v", err)
		return
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Log(response.ErrorLevel, "Database connection failed: %v", err)
		return
	}
	defer db.Close()

	router := gin.Default()

	health := db.Health()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, health)
	})

	usersDB := &users.Users{
		MainDB:    db.Main,
		ReplicaDB: db.Replica,
	}

	users.UsersSetUp(usersDB, router)

	port := config.LoadEnv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Log(response.InfoLevel, "Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Log(response.ErrorLevel, "Failed to start server: %v", err)
	}
}
