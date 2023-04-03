package exception

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/learn-micorservices/auth-service/model/web"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	response := web.WebResponse{
		Code:    code,
		Status:  false,
		Message: err.Error(),
		Data:    err.Error(),
	}

	return ctx.Status(code).JSON(response)
}
