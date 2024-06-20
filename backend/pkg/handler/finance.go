package handler

import (
	"strconv"
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addSchet(c *gin.Context) {
	req := &model.ReqSchet{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if len(req.ShopIDs) == 0 {
		bindShop, err := getBindShop(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		req.ShopIDs = []int{bindShop}
	}
	schet, err := h.services.AddSchet(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(schet, c)
}

func (h *Handler) getAllSchet(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetAllSchet(&filter)
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

func (h *Handler) getSchet(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, err := h.services.GetSchet(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateSchet(c *gin.Context) {
	req := model.ReqSchet{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	schet := &model.Schet{
		ID:           req.ID,
		Name:         req.Name,
		Currency:     req.Currency,
		Type:         req.Type,
		StartBalance: req.StartBalance,
		Deleted:      false,
	}
	err = h.services.UpdateSchet(schet)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteSchet(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteSchet(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}
