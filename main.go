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

	profileRepository := repository.NewProfileRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	profileService := service.NewProfileService(profileRepository, roleRepository, validate)
	profileController := controller.NewProfileController(profileService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	profileController.NewProfileRouter(app)

	err := app.Listen(serverConfig.URI)
	log.Println(err)
}

func main() {
	time.Local = time.UTC
	controllers()
}
