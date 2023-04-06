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
	user.Post("/check-token", controller.CheckToken)
}

func (controller *authController) Register(ctx *fiber.Ctx) error {
	request := new(web.RegisterRequest)
	err := ctx.BodyParser(request)
	helper.PanicIfError(err)

	user, err := controller.AuthService.Register(ctx.Context(), *request)
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

func (controller *authController) Login(ctx *fiber.Ctx) error {
	request := new(web.LoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	cookie, user, err := controller.AuthService.Login(ctx.Context(), *request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	ctx.Cookie(&cookie)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

func (controller *authController) Logout(ctx *fiber.Ctx) error {
	cookie := controller.AuthService.Logout(ctx.Context())

	ctx.Cookie(&cookie)

	// produce to kafka system log
	// cookieData := ctx.Cookies("token")
	// actor, _, _, _ := helper.ParseJwt(cookieData)
	// action := fmt.Sprintf("logout user %s", actor)
	// data := web.LogCreateRequest{
	// 	Actor:     actor,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(data, "POST.LOG", helper.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
	})
}

func (controller *authController) ForgetPassword(ctx *fiber.Ctx) error {
	data := new(web.ForgetPasswordRequest)
	err := ctx.BodyParser(data)
	helper.PanicIfError(err)

	if err := controller.AuthService.ForgetPassword(ctx.Context(), data.Email); err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	// //produce to kafka system log
	// action := fmt.Sprintf("request reset password %s", data.Email)
	// cookieData := ctx.Cookies("token")
	// actor, _, _, _ := helper.ParseJwt(cookieData)
	// dataLog := web.LogCreateRequest{
	// 	Actor:     actor,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(dataLog, "POST.LOG", helper.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset password has been sent",
	})
}

func (controller *authController) ResetPassword(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	token := ctx.Query("token")

	if len(token) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "token is missing",
		})
	}
	if len(email) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "email is missing",
		})
	}

	var request web.ResetPassword
	err := ctx.BodyParser(&request)
	helper.PanicIfError(err)

	controller.AuthService.ResetPassword(ctx.Context(), email, token, request)

	//produce to kafka system log
	// cookieData := ctx.Cookies("token")
	// actor, _, _, _ := helper.ParseJwt(cookieData)
	// action := fmt.Sprintf("reset password %s", email)
	// data := web.LogCreateRequest{
	// 	Actor:     actor,
	// 	Action:    action,
	// 	Timestamp: time.Now(),
	// }
	// helper.ProduceToKafka(data, "POST.LOG", helper.KafkaLogTopic)

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Reset successfully",
	})
}

func (controller *authController) CheckToken(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	token := ctx.Query("token")

	err := controller.AuthService.CheckToken(ctx.Context(), token, email)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    nil,
	})
}
