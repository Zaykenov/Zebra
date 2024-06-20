package handler

import (
	"strconv"
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addDealer(c *gin.Context) {
	req := &model.Dealer{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.Dealer.AddDealer(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getAllDealer(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Dealer.GetAllDealer(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if filter.Page == 0 {
		sendGeneral(res, c)
	} else {
		sendPagination(filter.Page, pageCount, res, c)
	}
}

func (h *Handler) getDealer(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetDealer(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateDealer(c *gin.Context) {
	req := model.Dealer{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.Dealer.UpdateDealer(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteDealer(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.Dealer.DeleteDealer(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}
