package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorMessage(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
}

func BadRequestError(ctx *fiber.Ctx, msg string, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message": msg,
		"error": err.Error(),
	})
}

func SuccessResponse(ctx *fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": message,
		"data":    data,
	})
}
