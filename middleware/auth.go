package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/learn-micorservices/auth-service/helper"
	"gitlab.com/learn-micorservices/auth-service/model/web"
)

func IsAuthenticated(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("token")

	if cookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	claims, err := helper.ParseJWT(cookie)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return ctx.Status(401).JSON(web.WebResponse{
				Code:    99281,
				Status:  false,
				Message: "token expired",
			})
		}

		return fiber.NewError(fiber.StatusUnauthorized)
	}

	ctx.Locals("claims", claims)

	return ctx.Next()
}

func IsAdmin(ctx *fiber.Ctx) error {
	claims := ctx.Locals("claims").(helper.JWTClaims)

	if claims.Profile.RoleID != "1" {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return ctx.Next()
}