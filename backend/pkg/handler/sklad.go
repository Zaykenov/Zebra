package handler

import (
	"errors"
	"strconv"
	"time"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RecalculateOneInventarization(c *gin.Context) {
	err := h.services.RecalculateInventarization()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) RecalculateInventarization(c *gin.Context) {
	err := h.services.RecalculateInventarization()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) RecalculateNetCost(c *gin.Context) {
	err := h.services.RecalculateNetCost()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) addSklad(c *gin.Context) {
	req := &model.Sklad{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.Sklad.AddSklad(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getAllSklad(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Sklad.GetAllSklad(&filter)
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

func (h *Handler) getSklad(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetSklad(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateSklad(c *gin.Context) {
	req := model.Sklad{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateSklad(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteSklad(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteSklad(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}
func (h *Handler) addToSkladWorker(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	// role, err := getUserRole(c)
	// if err != nil {
	// 	h.defaultErrorHandler(c, err)
	// 	return
	// }
	bindShop, err := getBindShop(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sklad, err := h.services.GetSkladByShopID(bindShop)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req := model.Postavka{}

	err = req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.SkladID = sklad.ID
	res, err := h.services.Transaction.GetLastShift(bindShop)

	if err != nil {
		if err.Error() == errors.New("record not found").Error() {
			h.defaultErrorHandler(c, errors.New("open shift"))
			return
		}
		h.defaultErrorHandler(c, err)
		return
	}
	if res.IsClosed {
		h.defaultErrorHandler(c, errors.New("open shift"))
		return
	}
	_, err = h.services.AddToSklad(&req, bindShop, id)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	// transaction := model.Transaction{
	// 	WorkerID: id,
	// 	SchetID:  req.SchetID,
	// 	ShopID:   bindShop,
	// 	Category: utils.Postavka,
	// 	Time:     req.Time,
	// 	Sum:      req.Sum,
	// 	Comment:  "Поставка №" + strconv.Itoa(postavka.ID),
	// }
	// err = h.services.PostavkaTransaction(role, &transaction, postavka.ID)
	// if err != nil {
	// 	h.defaultErrorHandler(c, err)
	// 	return
	// }
	sendSuccess(c)
}

func (h *Handler) addToSklad(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	// role, err := getUserRole(c)
	// if err != nil {
	// 	h.h.defaultErrorHandler(c, err)
	// 	return
	// }
	req := model.Postavka{}
	err = req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sklad, err := h.services.GetSklad(req.SkladID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	_, err = h.services.AddToSklad(&req, sklad.ShopID, id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	// transaction := model.Transaction{
	// 	WorkerID: id,
	// 	SchetID:  req.SchetID,
	// 	ShopID:   sklad.ShopID,
	// 	Category: utils.Postavka,
	// 	Time:     req.Time,
	// 	Sum:      req.Sum,
	// 	Comment:  "Поставка №" + strconv.Itoa(postavka.ID),
	// }
	// err = h.services.PostavkaTransaction(role, &transaction, postavka.ID)
	// if err != nil {
	// 	h.defaultErrorHandler(c, err)
	// 	return
	// }
	sendSuccess(c)
}

func (h *Handler) ostatki(c *gin.Context) {

	filter := model.Filter{}

	//items, err := h.services.Ostatki()
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Sklad.Ostatki(&filter)

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

func (h *Handler) RemoveFromSklad(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req := model.RemoveFromSklad{}

	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.WorkerID = userId

	if role == utils.ManagerRole || role == utils.MasterRole {
		//if manager
		err = h.services.RemoveFromSklad(&req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	} else if role == utils.WorkerRole {
		//else if worker then..
		bindShop, err := getBindShop(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		sklad, err := h.services.GetSkladByShopID(bindShop)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		req.SkladID = sklad.ID
		err = h.services.RequestToRemove(&req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	} else {
		h.defaultErrorHandler(c, errors.New("invalid role"))
		return
	}

	sendSuccess(c)
}

func (h *Handler) ConfirmToRemove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if role == utils.ManagerRole {
		err = h.services.ConfirmToRemove(id)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	} else {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) RejectToRemove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if role == utils.ManagerRole {
		err = h.services.RejectToRemove(id)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
	} else {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) UpdateSpisanie(c *gin.Context) {
	//Get role

	userId, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	req := model.RemoveFromSklad{}

	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.WorkerID = userId

	err = h.services.UpdateSpisanie(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) DeleteSpisanie(c *gin.Context) {
	req := model.ReqID{}

	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err := h.services.DeleteSpisanie(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) GetRemovedByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Sklad.GetRemovedByID(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)

}

func (h *Handler) GetRemoved(c *gin.Context) {
	/*sort := c.Request.URL.Query().Get("sortBy")
	if sort == "" {
		sort = "id asc"
	}
	if sort != "id asc" && sort != "id desc" && sort != "sklad_id asc" && sort != "sklad_id desc" && sort != "cost asc" && sort != "cost desc" {
		sort = "id asc"
	}
	search := c.Request.URL.Query().Get("search")

	sklad := c.Request.URL.Query().Get("sklad")
	res, err := h.services.GetRemoved(sort, sklad, search)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}*/

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Sklad.GetRemoved(&filter)

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

func (h *Handler) getItems(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetItems(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) GetAllPostavka(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.GetAllPostavka(&filter)
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

func (h *Handler) GetPostavka(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetPostavka(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) UpdatePostavka(c *gin.Context) {
	req := model.Postavka{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdatePostavka(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) DeletePostavka(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeletePostavka(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) AddTransfer(c *gin.Context) {
	req := &model.Transfer{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.Worker = id
	err = h.services.Sklad.AddTransfer(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) GetTransfer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetTransfer(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) GetAllTransfer(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.GetAllTransfer(&filter)
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

func (h *Handler) UpdateTransfer(c *gin.Context) {
	req := model.Transfer{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	req.Worker = id

	err = h.services.UpdateTransfer(&req)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)

}

func (h *Handler) DeleteTransfer(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteTransfer(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) CreateInventarization(c *gin.Context) {
	shops, err := getAvailableShops(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	req := model.Inventarization{}

	err = req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	hasAccess := false
	for _, shop := range shops {
		sklad, err := h.services.GetSkladByShopID(shop)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		if sklad.ID == req.SkladID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		h.defaultErrorHandler(c, errors.New("access restricted"))
		return
	}

	res, err := h.services.GetToCreateInventratization(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) GetInventarization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetInventarization(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) GetAllInventarization(c *gin.Context) {
	req := model.Filter{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.GetAllInventarization(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if req.Page == 0 {
		sendGeneral(res, c)
	} else {
		sendPagination(req.Page, pageCount, res, c)
	}
}

func (h *Handler) UpdateInventarization(c *gin.Context) {
	req := model.Inventarization{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if req.Status == utils.StatusOpened {
		res, err := h.services.OpenInventarization(&req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		sendGeneral(res, c)
	} else {
		res, err := h.services.UpdateInventarizationV2(&req)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		sendGeneral(res, c)
	}
}

func (h *Handler) UpdateInventarizationParams(c *gin.Context) {
	req := model.Inventarization{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.UpdateInventarizationParams(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) DeleteInventarization(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteInventarization(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) DeleteInventarizationItem(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteInventarizationItem(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) GetInventarizationDetailsIncome(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetInventarizationDetailsIncome(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) GetInventarizationDetailsExpence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetInventarizationDetailsExpence(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) GetInventarizationDetailsSpisanie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetInventarizationDetailsSpisanie(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) GetTrafficReport(c *gin.Context) {
	req := model.Filter{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if req.To.Day() == (time.Now().AddDate(0, 0, 2).Day()) {
		req.To = req.From
	}

	res, pageCount, err := h.services.GetTrafficReport(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if req.Page == 0 {
		sendGeneral(res, c)
	} else {
		sendPagination(req.Page, pageCount, res, c)
	}
}

func (h *Handler) GetAllInventarizationGroup(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetAllInventarizationGroup(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) CreateInventarizationGroup(c *gin.Context) {
	req := &model.ReqInventarizationGroup{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err := h.services.CreateInventarizationGroup(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) GetInventarizationGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetInventarizationGroup(filter, id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) DeleteInventarizationGroupItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.DeleteInventarizationGroup(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) UpdateInventarizationGroup(c *gin.Context) {
	req := &model.ReqInventarizationGroup{}
	if err := req.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err := h.services.UpdateInventarizationGroup(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) RecalculateDailyStatistic(c *gin.Context) {
	go h.services.RecalculateDailyStatistic()

	sendSuccess(c)
}
