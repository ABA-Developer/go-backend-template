package http

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/core/validator"
	"be-dashboard-nba/internal/modules/auth"
	"be-dashboard-nba/internal/modules/role/domain"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateRole godoc
// @Summary      Create Role
// @Description  Creates a new role. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role body      CreateRoleRequest true "Role data to create"
// @Success      201  {object}  presenter.ResponsePayloadData   "Successfully create role"
// @Failure      400  {object}  presenter.ResponseErrorPayload  "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed create role"
// @Security     BearerAuth
// @Router       /roles [post]
func CreateRole(usecase domain.RoleUsecase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request CreateRoleRequest
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

		err = usecase.CreateRoleUsecase(c.UserContext(), request.ToDomain(ah.GetClaims().UserID))
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

// UpdateRole godoc
// @Summary      Update Role
// @Description  Updates an existing role by ID. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id   path      int  true "Role ID"
// @Param        role body      UpdateRoleRequest true "Role data to update"
// @Success      200  {object}  presenter.ResponsePayloadData   "Successfully update role"
// @Failure      400  {object}  presenter.ResponseErrorPayload  "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed update role"
// @Security     BearerAuth
// @Router       /roles/{id} [put]
func UpdateRole(usecase domain.RoleUsecase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request UpdateRoleRequest
		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}
		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}
		roleID, err := c.ParamsInt("role_id")
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid role id",
			})
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		err = usecase.UpdateRoleUsecase(c.UserContext(), request.ToDomain(ah.GetClaims().UserID, roleID))
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed update role",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update role",
		})
	}
}

// ReadRoles godoc
// @Summary      Get Roles List
// @Description  Retrieves a paginated and searchable list of all roles. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        search query    string false "Search term for role name"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  presenter.ResponsePayloadPaginate{data=[]ReadRoleResponse} "Successfully get roles"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid query parameters"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get role"
// @Security     BearerAuth
// @Router       /roles [get]
func ReadRoles(usecase domain.RoleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request ReadRolesRequest
		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid query parameters",
			})
		}

		data, count, err := usecase.ReadRoleUsecase(c.UserContext(), request.ToDomainFilter())
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get role",
			})
		}

		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:    http.StatusOK,
			Message: "Successfully get roles",
			Data:    ToReadRoleResponses(data),
			Pagination: presenter.Pagination{
				Page:       request.Page,
				PageSize:   request.Limit,
				TotalItems: count, // A proper implementation determines TotalPages locally, omit for brevity mapped exactly to UI
			},
		})
	}
}

// ReadRoleDetail godoc
// @Summary      Get Role Detail
// @Description  Retrieves role details by ID. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        id   path      int  true "Role ID"
// @Success      200  {object}  presenter.ResponsePayloadData{data=ReadRoleResponse} "Successfully get role detail"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid role ID parameter"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404  {object}  presenter.ResponsePayloadMessage "Role not found"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get role detail"
// @Security     BearerAuth
// @Router       /roles/{id} [get]
func ReadRoleDetail(usecase domain.RoleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		roleID, err := c.ParamsInt("role_id")
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid role id",
			})
		}
		data, err := usecase.ReadDetailRoleUsecase(c.UserContext(), roleID)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get role detail",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    ToReadRoleResponse(data),
			Message: "Successfully get role detail",
		})
	}
}

// DeleteRole godoc
// @Summary      Delete Role
// @Description  Deletes an existing role by ID. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  presenter.ResponsePayloadData   "Successfully delete role"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid role ID parameter"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed delete role"
// @Security     BearerAuth
// @Router       /roles/{id} [delete]
func DeleteRole(usecase domain.RoleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		roleID, err := c.ParamsInt("role_id")
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid role id",
			})
		}
		err = usecase.DeleteRoleUsecase(c.UserContext(), roleID)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed delete role",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete role",
		})
	}
}

// ReadRoleAccess godoc
// @Summary      Get Role accesses Menu Permissions
// @Description  Retrieves a paginated and searchable list of menu permissions for a specific role. Requires Bearer token.
// @Tags         RoleAccess
// @Produce      json
// @Param        id     path     int    true "Role ID"
// @Param        search query    string false "Search term for menu name"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  presenter.ResponsePayloadPaginate{data=[]RoleAccessResponse} "Successfully get role detail"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid query parameters"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get get role detail"
// @Security     BearerAuth
// @Router       /roles/{id}/authorizations [get]
func ReadRoleAccess(usecase domain.RoleUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request ReadRoleAccessesRequest
		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid query parameters",
			})
		}
		roleIDStr := c.Params("role_id")

		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid role id",
			})
		}

		data, count, err := usecase.ReadRoleAccessUsecase(c.UserContext(), request.ToDomainFilter(roleID))
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get role accesses",
			})
		}
		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:    http.StatusOK,
			Message: "Successfully get role accesses",
			Data:    ToReadRoleAccessResponse(data),
			Pagination: presenter.Pagination{
				Page:       request.Page,
				PageSize:   request.Limit,
				TotalItems: count,
			},
		})
	}
}

// UpdateRoleAccess godoc
// @Summary      Update Role Access
// @Description  Updates access for multiple menu permissions for a specific role. Requires Bearer token.
// @Tags         RoleAccess
// @Accept       json
// @Produce      json
// @Param        id         path     int                       true "Role ID"
// @Param        access     body     UpdateRoleAccessRequest   true "List of role access updates"
// @Success      200  {object}  presenter.ResponsePayloadData   "Successfully update role accesses"
// @Failure      400  {object}  presenter.ResponseErrorPayload  "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed update role accesses"
// @Security     BearerAuth
// @Router       /roles/{id}/authorizations [patch]
func UpdateRoleAccess(usecase domain.RoleUsecase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request UpdateRoleAccessRequest
		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}
		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}
		roleIDStr := c.Params("role_id")
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Invalid role id",
			})
		}

		err = usecase.UpdateRoleAccessUsecase(c.UserContext(), roleID, request.ToDomainPayloads(roleID))
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed update role accesses",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update role accesses",
		})
	}
}
