package handler

import (
	"errors"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Registrate(c *gin.Context) {
	var reqRegistrate model.ReqRegistrate
	err := reqRegistrate.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	response, err := h.services.Registrate(&reqRegistrate)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if response == utils.AlreadyRegistered {
		h.defaultErrorHandler(c, errors.New("this email already used"))
		return
	}
	sendGeneral("codeSentToEmail", c)
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var reqVerify model.ReqVerifyEmail
	err := reqVerify.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	id, response, err := h.services.VerifyEmail(&reqVerify)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if response == utils.IncorrectCode {
		h.defaultErrorHandler(c, errors.New("incorrect code"))
		return
	}
	sendGeneral(id, c)
}

func (h *Handler) signInClient(c *gin.Context) {
	var reqSignIn model.ReqSignIn
	err := reqSignIn.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.SendCode(reqSignIn.Email)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral("codeSentToEmail", c)
}
