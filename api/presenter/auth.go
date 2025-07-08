package presenter

import (
	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func AuthSuccessResponse(data *Auth) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

func AuthErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
