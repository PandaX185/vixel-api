package main

import (
	"log"
	"net/http"
	"vixel/config"
	"vixel/domains/image"
	"vixel/domains/processing"
	"vixel/domains/user"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	if err := config.LoadEnvConfig(); err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}
	app := gin.Default()
	c := cors.AllowAll()
	app.Use(func(gc *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gc.Next()
		})
		c.Handler(next).ServeHTTP(gc.Writer, gc.Request)
	})
	api := app.Group("/api/v1")

	db, err := config.NewPostgres()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.Migrator().AutoMigrate(&user.User{}, &image.Image{})

	userService := user.NewUserService(db)
	userHandler := user.NewUserHandler(userService)
	userHandler.SetupUserRoutes(api)

	imageService := image.NewImageService(db)
	uploadService := image.NewUploadService()
	imageHandler := image.NewImageHandler(imageService, uploadService)
	imageHandler.SetupImageRoutes(api)

	processingService := processing.NewProcessingService(db, uploadService)
	processingHandler := processing.NewProcessingHandler(processingService)
	processingHandler.SetupProcessingRoutes(api)

	if err := app.Run(config.Config.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
