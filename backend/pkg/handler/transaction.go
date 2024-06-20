package handler

import (
	"errors"
	"strconv"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createTransaction(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	req := model.Transaction{}

	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if role == utils.WorkerRole {
		if req.SchetID == 0 {
			shop, err := h.services.Shop.GetShopByID(bindShop)
			if err != nil {
				h.defaultErrorHandler(c, err)
				return
			}
			req.SchetID = shop.CashSchet
			if shop.Blocked {
				h.defaultErrorHandler(c, errors.New("shop is not available"))
			}
		}
	}
	req.WorkerID = id
	switch req.Category {
	case utils.OpenShift:
		if role == utils.ManagerRole {
			h.defaultErrorHandler(c, errors.New("bad request | manager can't open shift"))
			return
		}
		err = h.services.Transaction.OpenShift(id, &req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	case utils.CloseShift:
		if role == utils.ManagerRole && req.ShiftID == 0 {
			h.defaultErrorHandler(c, errors.New("bad request | can't close shift without id"))
			return
		}
		err = h.services.Transaction.CloseShift(id, &req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	case utils.Collection:
		err := h.services.Transaction.Collection(role, id, &req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	case utils.Income:
		err := h.services.Transaction.Income(role, id, &req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	default:
		h.defaultErrorHandler(c, errors.New("bad request | category is invalid"))
	}
	sendSuccess(c)
}

func (h *Handler) getAllShifts(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.Transaction.GetAllShifts(&filter)
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

func (h *Handler) getShiftByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Transaction.GetShiftByID(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getShiftByShopId(c *gin.Context) {
	workerID, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	id, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	hasShift, err := h.services.Transaction.GetShiftByShopId(id, workerID)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(hasShift, c)
}

func (h *Handler) getAllTransaction(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Transaction.GetAllTransaction(&filter)
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

func (h *Handler) getTransactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Transaction.GetTransactionByID(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateTransaction(c *gin.Context) {
	workerID, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req := model.Transaction{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	blocked, err := h.services.Transaction.CheckForBlockedShop(req.ShiftID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if blocked {
		h.defaultErrorHandler(c, errors.New("shop is not available"))
	}
	req.UpdatedWorkerID = workerID
	err = h.services.Transaction.UpdateTransaction(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) deleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.Transaction.DeleteTransaction(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}
