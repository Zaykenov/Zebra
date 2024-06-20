package handler

import (
	"errors"
	"strconv"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllWorkers(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	workers, pageCount, err := h.services.GetAllWorkers(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if filter.Page == 0 {
		sendGeneral(workers, c)
	} else {
		sendPagination(filter.Page, pageCount, workers, c)
	}

}

func (h *Handler) getWorker(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	worker, err := h.services.GetWorker(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(worker, c)
}

func (h *Handler) updateWorker(c *gin.Context) {
	var reqWorker model.Worker
	err := reqWorker.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if reqWorker.Role == utils.WorkerRole && len(reqWorker.Shops) > 1 {
		h.defaultErrorHandler(c, errors.New("worker can have only one shop"))
		return
	}
	worker, err := h.services.UpdateWorker(&reqWorker)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(worker, c)
}

func (h *Handler) deleteWorker(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.DeleteWorker(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) Terminal(c *gin.Context) {
	// bindShop, err := getBindShop(c)
	// if err != nil {
	// 	h.defaultErrorHandler(c, err)
	// 	return
	// }
	// res, err := h.services.GetAllTerminal
}

func (h *Handler) PingPong(c *gin.Context) {
	sendGeneral("Pong!", c)
}
