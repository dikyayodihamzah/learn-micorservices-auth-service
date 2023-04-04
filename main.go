package main

import (
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/learn-micorservices/auth-service/config"
	"gitlab.com/learn-micorservices/auth-service/controller"
	"gitlab.com/learn-micorservices/auth-service/repository"
	"gitlab.com/learn-micorservices/auth-service/service"
)

func controllers() {
	time.Local = time.UTC

	serverConfig := config.NewServerConfig()
	db := config.NewDB
	validate := validator.New()

	authRepository := repository.NewAuthRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	authService := service.NewAuthService(authRepository, roleRepository, validate)
	authController := controller.NewAuthController(authService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	authController.NewAuthRouter(app)

	err := app.Listen(serverConfig.URI)
	log.Println(err)
}

func main() {
	time.Local = time.UTC
	controllers()
}
