package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/learn-micorservices/auth-service/config"
	"gitlab.com/learn-micorservices/auth-service/exception"
	"gitlab.com/learn-micorservices/auth-service/helper"
	"gitlab.com/learn-micorservices/auth-service/middleware"
	"gitlab.com/learn-micorservices/auth-service/model/web"
	"gitlab.com/learn-micorservices/auth-service/service"
)

type AuthController interface {
	NewAuthRouter(app *fiber.App)
}
type authController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		AuthService: authService,
	}
}

func (controller *authController) NewAuthRouter(app *fiber.App) {
	user := app.Group(config.EndpointPrefix)

	user.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
			Code:    fiber.StatusOK,
			Status:  true,
			Message: "ok",
		})
	})

	// user.Use(middleware.IsAuthenticated)
	user.Post("/register", controller.Register)
	user.Post("/login", controller.Login)
	user.Post("/logout", middleware.IsAuthenticated, controller.Logout)

	user.Post("/forget-password", controller.ForgetPassword)
	user.Post("/reset-password", controller.ResetPassword)
	user.Post("/refresh", controller.Refersh)
	user.Post("/check-token", controller.CheckToken)
}


func (controller *authController) Register(ctx *fiber.Ctx) error {

	user, err := controller.AuthService.GetCurrentProfile(ctx.Context())
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

func (controller *authController) UpdateProfile(ctx *fiber.Ctx) error {
	claims := ctx.Locals("claims").(helper.JWTClaims)

	request := new(web.UpdateProfileRequest)
	if err := ctx.BodyParser(request); err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	user, err := controller.ProfileService.UpdateProfile(ctx.Context(), claims, *request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}