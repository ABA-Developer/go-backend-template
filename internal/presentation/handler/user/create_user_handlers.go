package handlers

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/validator"
	"be-dashboard-nba/internal/presentation/presenter"
	userPresenter "be-dashboard-nba/internal/presentation/presenter/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
)

// Summary      Create User
// @Description  Creates a new user with the provided details. Requires Bearer token.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body      userPresenter.CreateUserRequest true "User data to create"
// @Success      201  {object}  presenter.ResponsePayloadData     "Successfully create user"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404  {object}  presenter.ResponsePayloadMessage "Role ID not found"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /users [post]
func CreateUser(svc contract.UserUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.CreateUserRequest
		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.CreateUserUseCase(c.UserContext(), request, userID)
		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrRoleIdNotFound.Message,
				})
			}
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed create user",
				})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create user",
		})

	}
}
