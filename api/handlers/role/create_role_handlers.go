package handlers

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/role/service"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateRole godoc
// @Summary      Create Role
// @Description  Creates a new role. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role body      presenter.CreateRoleRequest true "Role data to create"
// @Success      201  {object}  presenter.ResponsePayloadData   "Successfully create role"
// @Failure      400  {object}  presenter.ResponseErrorPayload  "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed create role"
// @Security     BearerAuth
// @Router       /roles [post]
func CreateRole(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		var payload rolePresenter.CreateRoleRequest
		if err = c.BodyParser(&payload); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(payload); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		fmt.Printf("userId: %s", userID)

		err = svc.CreateRoleService(c.UserContext(), payload, userID)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed create role",
				})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create role",
		})
	}
}
