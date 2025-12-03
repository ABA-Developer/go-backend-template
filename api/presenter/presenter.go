package presenter

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ResponsePayloadMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseMessage(c *fiber.Ctx, res ResponsePayloadMessage) error {
	return c.Status(res.Code).JSON(res)
}

type ResponsePayloadData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseData(c *fiber.Ctx, res ResponsePayloadData) error {
	return c.Status(res.Code).JSON(res)
}

type Pagination struct {
	Page       int  `json:"page" example:"1"`
	PageSize   int  `json:"page_size" example:"20"`
	TotalItems int  `json:"total_items" example:"100"`
	TotalPages int  `json:"total_pages" example:"5"`
	HasNext    bool `json:"has_next" example:"true"`
	HasPrev    bool `json:"has_prev" example:"false"`
}

type ResponsePayloadPaginate struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func ResponsePaginate(c *fiber.Ctx, res ResponsePayloadPaginate) error {
	return c.Status(res.Code).JSON(res)
}

type ResponseErrorPayload struct {
	Code    int         `json:"code"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message"`
}

func ResponseErrorValidate(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(ResponseErrorPayload{
		Code:    http.StatusBadRequest,
		Error:   validator.ValidationErrors(err),
		Message: constant.ErrMsgValidate,
	})
}
