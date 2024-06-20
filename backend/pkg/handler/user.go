package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Get all workers
// @Tags worker
// @Description get all workers
// @ID get-all-workers
// @Accept  json
// @Produce  json
// @Param input body model.ReqLogin false "login"
// @Success 200 {object} model.Worker
// @Router /signin [post]
func (h *Handler) signIn(c *gin.Context) {
	var reqLogin model.ReqLogin
	err := reqLogin.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	worker, err := h.services.GetWorkerByUsername(reqLogin.Username)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if worker.Username != reqLogin.Username || worker.Password != reqLogin.Password {
		h.defaultErrorHandler(c, errors.New("wrong username or password"))
		return
	}
	sendGeneral(worker, c)
}

func (h *Handler) authorize(c *gin.Context) {
	var reqUser model.ReqWorkerRegistration
	err := reqUser.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if reqUser.Role == utils.WorkerRole && len(reqUser.Shops) > 1 {
		h.defaultErrorHandler(c, errors.New("worker can have only one shop"))
		return
	}
	worker, err := h.services.GetWorkerByUsername(reqUser.Username)
	if err != nil {
		if err.Error() != "record not found" {
			h.defaultErrorHandler(c, err)
			return
		}
	}

	if worker != nil {
		h.defaultErrorHandler(c, errors.New("user with this username already exists. Please, try another username"))
		return
	}

	err = h.services.CreateWorker(&reqUser)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetWorkerByUsername(reqUser.Username)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) addUser(c *gin.Context) {
	input := &model.UserRequest{}
	if err := input.ParseRequest(c); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.User.AddUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.services.User.GetUser(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getCurrentOrders(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orders, err := h.services.User.GetCurrentOrders(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) addFeedback(c *gin.Context) {
	input := &model.FeedbackRequest{}
	if err := input.ParseRequest(c); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.User.AddFeedback(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) genCode(c *gin.Context) { //?
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	code, timeDate, err := h.services.User.GenCode(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	tm := time.Unix(timeDate, 0)
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":  code,
		"until": tm,
	})
}

func (h *Handler) getUserByCode(c *gin.Context) {
	code := c.Param("code")
	user, err := h.services.User.GetUserByCode(code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getAllMobileUsers(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	log.Println(filter)
	users, pageCount, err := h.services.GetAllMobileUsers(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendPagination(filter.Page, pageCount, users, c)
}

func (h *Handler) getMobileUser(c *gin.Context) {
	availableShops, err := getAvailableShops(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	id := c.Param("id")
	user, err := h.services.GetMobileUser(id, availableShops)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(user, c)
}

func (h *Handler) getAllFeedbacks(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	feedbacks, pageCount, err := h.services.GetAllFeedbacks(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendPagination(filter.Page, pageCount, feedbacks, c)
}

func (h *Handler) getFeedback(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	feedback, err := h.services.GetFeedback(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(feedback, c)
}
