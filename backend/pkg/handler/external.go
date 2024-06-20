package handler

import (
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) saveCheck(c *gin.Context) {
	req := &model.ReqTisResponse{}
	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.SaveCheck(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}
