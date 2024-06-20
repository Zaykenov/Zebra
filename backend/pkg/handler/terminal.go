package handler

import (
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) startTerminal(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	workerID, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetAllProducts(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	shift, err := h.services.GetShiftByShopId(filter.BindShop, workerID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res.CurrentShift = shift
	sendGeneral(res, c)
}
