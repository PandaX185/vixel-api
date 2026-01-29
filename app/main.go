package main

import (
	"log"
	"vixel/config"
	"vixel/domains/user"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadEnvConfig(); err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}
	app := gin.Default()
	api := app.Group("/api/v1")

	db, err := config.NewPostgres()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.Migrator().AutoMigrate(&user.User{})

	userService := user.NewUserService(db)
	userHandler := user.NewUserHandler(userService)
	userHandler.SetupUserRoutes(api)

	if err := app.Run(config.Config.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
