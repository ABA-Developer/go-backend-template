package http

import (
	presenter "be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ReadMenuDetail(usecase domain.MenuUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}

		data, err := usecase.ReadMenuDetailUsecase(c.UserContext(), menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    data,
			Message: "Successfully get menu detail",
		})

	}
}
