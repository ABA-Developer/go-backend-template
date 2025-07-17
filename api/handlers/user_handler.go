package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user"

)

func CreateUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.User
		if err := c.BodyParser(&requestBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		if requestBody.Email == "" || requestBody.Password == "" {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(errors.New(
				"Email and Password cannot be empty")))
		}

		result, err := service.CreateUser(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		return c.JSON(presenter.UserSuccessResponse(result))
	}
}

func GetAllUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetUsers()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(err))
		}
		return c.JSON(presenter.UsersSuccessResponse(result))
	}
}

func GetUserById(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		result, err := service.GetUserById(int64(id))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(err))
		}
		return c.JSON(presenter.UserSuccessResponse(result))
	}
}

func DeleteUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		err = service.DeleteUser(int64(id))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		result := &entities.User{
			ID: int64(id),
		}
		return c.JSON(presenter.UserSuccessResponse(result))
	}
}

func UpdateUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		var requestBody entities.User
		if err := c.BodyParser(&requestBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		requestBody.ID = int64(id)

		result, err := service.UpdateUser(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserErrorResponse(err))
		}

		return c.JSON(presenter.UserSuccessResponse(result))
	}
}
