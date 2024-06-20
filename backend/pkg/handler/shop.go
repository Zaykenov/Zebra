package handler

import (
	"net/http"
	"strconv"
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllShops(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, totalPages, err := h.services.Shop.GetAllShops(&filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if filter.Page == 0 {
		sendGeneral(shops, c)
	} else {
		sendPagination(filter.Page, totalPages, shops, c)
	}
}

func (h *Handler) getShopByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	shop, err := h.services.Shop.GetShopByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sendGeneral(shop, c)
}

func (h *Handler) createShop(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var reqShop model.ShopGlobal
	if err := reqShop.ParseRequest(c); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tech_carts, err := h.services.Master.GetAllTechCartsMasterIds()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	reqShop.ProductsShop.TechCarts = tech_carts

	shop, err := h.services.Shop.CreateShop(reqShop.Shop, reqShop.ProductsShop)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, worker := range reqShop.Workers {
		worker.BindShop = shop.ID
		worker.Shops = []int{shop.ID}
		if worker.New {
			err = h.services.CreateWorker(worker)
			if err != nil {
				h.defaultErrorHandler(c, err)
				return
			}
		} else {
			oldWorker, err := h.services.GetWorkerByUsername(worker.Username)
			if err != nil {
				h.defaultErrorHandler(c, err)
				return
			}
			if oldWorker != nil {
				exists := false
				for _, shopID := range oldWorker.Shops {
					if shopID == shop.ID {
						exists = true
						break
					}
				}
				if !exists {
					oldWorker.Shops = append(oldWorker.Shops, shop.ID)
				}
				_, err = h.services.UpdateWorker(oldWorker)
				if err != nil {
					h.defaultErrorHandler(c, err)
					return
				}
			}
		}
	}
	reqShop.Sklad.ShopID = shop.ID
	err = h.services.Sklad.AddSklad(reqShop.Sklad)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	worker, err := h.services.GetWorker(userID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	worker.Shops = append(worker.Shops, shop.ID)
	_, err = h.services.UpdateWorker(worker)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) updateShop(c *gin.Context) {
	var reqShop model.ReqShop
	if err := reqShop.ParseRequest(c); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Shop.UpdateShop(&reqShop)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(c)
}

func (h *Handler) deleteShop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Shop.DeleteShop(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccess(c)
}
