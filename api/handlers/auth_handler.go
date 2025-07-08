package handlers

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/pkg/auth"
	"be-dashboard-nba/pkg/entities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
		// var requestBody entities.Book
		// err := c.BodyParser(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusBadRequest)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// if requestBody.Author == "" || requestBody.Title == "" {
		// 	c.Status(http.StatusInternalServerError)
		// 	return c.JSON(presenter.BookErrorResponse(errors.New(
		// 		"Please specify title and author")))
		// }
		// result, err := service.InsertBook(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusInternalServerError)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// return c.JSON(presenter.BookSuccessResponse(result))
	}
}

func LoginUser(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Auth
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(presenter.AuthErrorResponse(err))
		}

		// Check for empty payload
		if requestBody.Email == "" || requestBody.Password == "" {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(presenter.AuthErrorResponse(errors.New(
				"Email and Password cannot be empty")))
		}

		// Chek for the email
		user, err := service.GetUserByEmail(requestBody.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(presenter.AuthErrorResponse(err))
		}

		// Compare encrypted password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(presenter.AuthErrorResponse(err))
		}

		// Create JWT
		claims := jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte("jwt_secret"))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(presenter.AuthErrorResponse(err))
		}

		res := &presenter.Auth{
			Email: user.Email,
			Token: t,
		}

		return c.JSON(presenter.AuthSuccessResponse(res))
	}
}

// UpdateBook is handler/controller which updates data of Books in the BookShop
func LogoutUser(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
		// var requestBody entities.Book
		// err := c.BodyParser(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusBadRequest)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// result, err := service.UpdateBook(&requestBody)
		// if err != nil {
		// 	c.Status(http.StatusInternalServerError)
		// 	return c.JSON(presenter.BookErrorResponse(err))
		// }
		// return c.JSON(presenter.BookSuccessResponse(result))
	}
}
