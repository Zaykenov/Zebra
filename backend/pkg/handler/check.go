package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addCheck(c *gin.Context) {

	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req := model.ReqCheck{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	req.WorkerID = id
	req.ShopID = bindShop
	res, err := h.services.Transaction.GetLastShift(bindShop)
	if err != nil {
		if err.Error() == errors.New("record not found").Error() {
			h.services.SaveError(errors.New("open shift").Error(), string(jsonReq))
			h.defaultErrorHandler(c, errors.New("open shift"))
			return
		}
		h.services.SaveError(err.Error(), string(jsonReq))
		h.defaultErrorHandler(c, err)
		return
	}

	if res.IsClosed {
		h.services.SaveError(errors.New("open shift").Error(), string(jsonReq))
		h.defaultErrorHandler(c, errors.New("open shift"))
		return
	}

	sklad, err := h.services.GetSkladByShopID(req.ShopID)
	if err != nil {
		h.services.SaveError(err.Error(), string(jsonReq))
		h.defaultErrorHandler(c, err)
		return
	}
	req.SkladID = sklad.ID
	oldCheck, err := h.services.GetCheckByIdempotency(req.IdempotencyKey)
	if err != nil {
		if err.Error() != errors.New("record not found").Error() {
			h.services.SaveError(err.Error(), string(jsonReq))
			h.defaultErrorHandler(c, err)
			return
		}
	}
	if oldCheck != nil {
		sendGeneral(oldCheck, c)
		return
	}
	if req.Status == utils.StatusOpened && req.ID == 0 {
		check, err := h.services.AddCheck(&req)
		if err != nil {
			h.services.SaveError(err.Error(), string(jsonReq))
			h.defaultErrorHandler(c, err)
			return
		}
		sendGeneral(check, c)
	} else {
		if req.ID != 0 {
			oldCheck, err := h.services.GetCheckByID(req.ID)
			if err != nil {
				h.services.SaveError(err.Error(), string(jsonReq))
				h.defaultErrorHandler(c, err)
				return
			}

			if oldCheck.Status == utils.StatusClosed {
				h.services.SaveError(errors.New("check is closed").Error(), string(jsonReq))
				h.defaultErrorHandler(c, errors.New("check is closed"))
				return
			}
		}
		check, err := h.services.UpdateCheck(&req)
		if err != nil {
			h.services.SaveError(err.Error(), string(jsonReq))
			h.defaultErrorHandler(c, err)
			return
		}
		sendGeneral(check, c)
	}
}

func (h *Handler) getAllCheckWorker(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetAllWorkerCheck(&filter)
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

func (h *Handler) getAllCheck(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetAllCheck(&filter)
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

func (h *Handler) DeleteCheck(c *gin.Context) {
	req := model.ReqID{}

	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	check, err := h.services.GetCheckByID(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if check.Status == utils.StatusClosed {
		h.defaultErrorHandler(c, errors.New("check is closed"))
		return
	}

	err = h.services.DeleteCheck(req.ID)

	if err != nil {
		h.defaultErrorHandler(c, err)
	}

	sendSuccess(c)
}

func (h *Handler) getAllWorkerCheck(c *gin.Context) {
	/*id, err := getUserId(c)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetAllWorkerCheck(id)

	if err != nil {
		defaultErrorHandler(c, err)
		return
	}*/

}

func (h *Handler) getCheckByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetCheckByID(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getCheckByIDForPrinter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetCheckByIDForPrinter(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) addTag(c *gin.Context) {
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req := model.Tag{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.ShopID = bindShop

	err = h.services.AddTag(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) getAllTag(c *gin.Context) {
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetAllTag(bindShop)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetTag(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) updateTag(c *gin.Context) {
	req := model.Tag{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err := h.services.UpdateTag(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) deleteTag(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.DeleteTag(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) DeactivateCheck(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.DeactivateCheck(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) IdempotencyCheck(c *gin.Context) {
	keys := &model.IdempotencyCheckArray{}
	if err := keys.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	var (
		bindShop int
		err      error
	)

	if keys.ShopID == 0 {
		bindShop, err = getBindShop(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	} else {
		bindShop = keys.ShopID
	}

	err = h.services.IdempotencyCheck(keys, bindShop)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) saveFailedCheck(c *gin.Context) {
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	req := []*model.FailedCheck{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	for _, r := range req {
		r.Created_at = time.Now()
		r.ShopID = bindShop
	}

	err = h.services.SaveFailedCheck(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) getStoliki(c *gin.Context) {
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetStoliki(bindShop)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getFilledStoliki(c *gin.Context) {
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetFilledStoliki(bindShop)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}
