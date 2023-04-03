package exception

import "github.com/gofiber/fiber/v2"

func ErrBadRequest(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func ErrUnauthorized(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnauthorized, message)
}

func ErrNotFound(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound, message)
}

func ErrUnprocessableEntity(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnprocessableEntity, message)
}

func ErrInternalServer(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
