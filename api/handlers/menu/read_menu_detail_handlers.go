package handlers

import (
	presenter "be-dashboard-nba/api/presenter"
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ReadMenuDetail(svc service.Service) fiber.Handler {
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

		data, err := svc.ReadMenuDetailService(c.UserContext(), menuID)
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
			Data:    menuPresenter.ToReadMenuDetailResponse(data),
			Message: "Successfully get menu detail",
		})

	}
}
