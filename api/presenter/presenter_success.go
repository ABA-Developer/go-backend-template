package presenter

import (
	"be-dashboard-nba/internal/utils"
	"math"

	"github.com/gofiber/fiber/v2"
)

type responseDataPayload struct {
	Code 	  int  		  `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	TotalData *int64      `json:"total_data,omitempty"`
}

type responsePaginatePayload struct {
	Data     interface{}     `json:"data"`
	Message  string          `json:"message"`
	Paginate paginatePayload `json:"paginate"`
}

type paginatePayload struct {
	CurrentPage int     `json:"current_page"`
	PerPage     int     `json:"per_page"`
	TotalPage   float64 `json:"total_page"`
	TotalData   int64   `json:"total_data"`
}

func ResponseData(c *fiber.Ctx, res ResponsePayload) error {
	return c.Status(res.Code).JSON(responseDataPayload{
		Code: res.Code,
		Data:    res.Data,
		Message: res.Message,
	})
}

func ResponsePaginate(c *fiber.Ctx, paginate utils.PaginationPayload, totalData int64, res ResponsePayload) error {
	if paginate.Page <= 0 && paginate.Limit <= 0 {
		return c.Status(res.Code).JSON(responseDataPayload{
			Data:      res.Data,
			Message:   res.Message,
			TotalData: &totalData,
		})
	}

	return c.Status(res.Code).JSON(responsePaginatePayload{
		Data:     res.Data,
		Message:  res.Message,
		Paginate: toPaginatePayload(paginate.Page, paginate.Limit, totalData),
	})
}

func toPaginatePayload(currentPage, perPage int, totalData int64) paginatePayload {
	totalPage := math.Ceil(float64(totalData) / float64(perPage))
	return paginatePayload{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPage:   totalPage,
		TotalData:   totalData,
	}
}
