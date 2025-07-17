package presenter

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/pkg/entities"
)

// Book is the presenter object which will be passed in the response by Handler
type User struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
}

func UserSuccessResponse(data *entities.User) *fiber.Map {
	user := User{
		ID:         data.ID,
		FirstName:  data.FirstName,
		MiddleName: data.MiddleName,
		LastName:   data.LastName,
		Email:      data.Email,
	}
	return &fiber.Map{
		"status": true,
		"data":   user,
		"error":  nil,
	}
}

func UsersSuccessResponse(datas []entities.User) *fiber.Map {
	var users []User
	for _, data := range datas {
		user := User{
			ID:         data.ID,
			FirstName:  data.FirstName,
			MiddleName: data.MiddleName,
			LastName:   data.LastName,
			Email:      data.Email,
		}
		users = append(users, user)
	}

	return &fiber.Map{
		"status": true,
		"data":   users,
		"error":  nil,
	}
}

func UserErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
