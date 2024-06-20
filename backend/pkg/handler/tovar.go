package handler

import (
	"strconv"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addTovar(c *gin.Context) {
	req := &model.ReqTovar{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	tovarMaster := &model.TovarMaster{
		Name:     req.Name,
		Category: req.Category,
		Image:    req.Image,
		Tax:      req.Tax,
		Measure:  req.Measure,
		Price:    req.Price,
		Discount: req.Discount,
		Deleted:  false,
	}

	if role == utils.MasterRole {
		tovarMaster.Status = utils.MenuStatusApproved
	} else {
		tovarMaster.Status = utils.MenuStatusPending
	}
	newTovar, err := h.services.Master.AddTovarMaster(tovarMaster)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if role == utils.MasterRole {
		shops, err := getAvailableShops(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		req.ShopID = shops
	}
	tovars := []*model.Tovar{}
	for _, shop := range req.ShopID {
		tovar := &model.Tovar{
			ShopID:    shop,
			TovarID:   newTovar.ID,
			Name:      newTovar.Name,
			Category:  newTovar.Category,
			Image:     newTovar.Image,
			Tax:       newTovar.Tax,
			Measure:   newTovar.Measure,
			Discount:  newTovar.Discount,
			IsVisible: true,
			Price:     newTovar.Price,
			Cost:      req.Cost,
		}
		tovars = append(tovars, tovar)
	}
	err = h.services.AddTovars(tovars)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getAllTovar(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Tovar.GetAllTovar(&filter)
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

func (h *Handler) getTovar(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Tovar.GetTovar(id, &filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateTovar(c *gin.Context) {
	req := model.ReqTovar{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateTovar(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteTovar(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteTovar(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) addCategoryTovar(c *gin.Context) {
	req := model.CategoryTovar{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.AddCategoryTovar(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getAllCategoryTovar(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetAllCategoryTovar(&filter)
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

func (h *Handler) updateCategoryTovar(c *gin.Context) {
	req := model.CategoryTovar{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateCategoryTovar(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteCategoryTovar(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteCategoryTovar(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getCategoryTovar(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetCategoryTovar(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) addTechCart(c *gin.Context) {
	req := model.ReqTechCart{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	ingredientsMaster := []*model.IngredientTechCartMaster{}
	for _, ingredient := range req.Ingredients {
		ingredientsMaster = append(ingredientsMaster, &model.IngredientTechCartMaster{
			TechCartID:   ingredient.TechCartID,
			IngredientID: ingredient.IngredientID,
			Brutto:       ingredient.Brutto,
		})
	}
	naborsMaster := []*model.NaborTechCartMaster{}
	for _, nabor := range req.Nabor {
		naborsMaster = append(naborsMaster, &model.NaborTechCartMaster{
			TechCartID: nabor.TechCartID,
			NaborID:    nabor.NaborID,
		})
	}
	techCartMaster := &model.TechCartMaster{
		Name:        req.Name,
		Category:    req.Category,
		Image:       req.Image,
		Price:       req.Price,
		Tax:         req.Tax,
		Measure:     req.Measure,
		Deleted:     false,
		Discount:    req.Discount,
		Ingredients: ingredientsMaster,
		Nabor:       naborsMaster,
	}
	if role == utils.MasterRole {
		techCartMaster.Status = utils.MenuStatusApproved
	} else {
		techCartMaster.Status = utils.MenuStatusPending
	}
	newTechCart, err := h.services.AddTechCartMaster(techCartMaster)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if role == utils.MasterRole {
		shops, err := getAvailableShops(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		req.ShopID = shops
	}

	techCarts := []*model.TechCart{}
	for _, shop := range req.ShopID {
		ingredients := []*model.IngredientTechCart{}
		for _, ingredient := range req.Ingredients {
			ingredients = append(ingredients, &model.IngredientTechCart{
				TechCartID:   newTechCart.ID,
				IngredientID: ingredient.IngredientID,
				Brutto:       ingredient.Brutto,
				ShopID:       shop,
			})
		}
		nabors := []*model.NaborTechCart{}
		for _, nabor := range req.Nabor {
			nabors = append(nabors, &model.NaborTechCart{
				TechCartID: newTechCart.ID,
				NaborID:    nabor.NaborID,
				ShopID:     shop,
			})
		}
		techCart := &model.TechCart{
			ShopID:      shop,
			TechCartID:  newTechCart.ID,
			Name:        newTechCart.Name,
			Category:    newTechCart.Category,
			Image:       newTechCart.Image,
			Tax:         newTechCart.Tax,
			Measure:     newTechCart.Measure,
			Discount:    newTechCart.Discount,
			Price:       newTechCart.Price,
			IsVisible:   true,
			Ingredients: ingredients,
			Nabor:       nabors,
		}
		techCarts = append(techCarts, techCart)
	}

	err = h.services.AddTechCarts(techCarts)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getAllTechCart(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetAllTechCart(&filter)
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

func (h *Handler) getTechCart(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetTechCart(id, &filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateTechCart(c *gin.Context) {
	req := model.ReqTechCart{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.UpdateTechCart(&req, role)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteTechCart(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteTechCart(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getTovarWithParams(c *gin.Context) {

	sort := c.Request.URL.Query().Get("sortBy")
	if sort == "" {
		sort = "name.asc"
	}
	category := c.Request.URL.Query().Get("category")
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		categoryID = 0
	}
	search := c.Request.URL.Query().Get("search")

	sklad := c.Request.URL.Query().Get("sklad")

	res, err := h.services.GetTovarWithParams(sort, sklad, search, categoryID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) getTechCartWithParams(c *gin.Context) {
	sort := c.Request.URL.Query().Get("sortBy")
	if sort == "" {
		sort = "name.asc"
	}
	category := c.Request.URL.Query().Get("category")
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sklad := c.Request.URL.Query().Get("sklad")

	res, err := h.services.GetTechCartWithParams(sort, sklad, categoryID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) getTechCartNabor(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetTechCartNabor(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) getEverything(c *gin.Context) {
	res, err := h.services.GetEverything()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getDeletedTovar(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetDeletedTovar(&filter)
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

func (h *Handler) getDeletedTechCart(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetDeletedTechCart(&filter)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	if filter.Page == 0 {
		sendGeneral(res, c)
	} else {
		sendPagination(filter.Page, pageCount, res, c)
	}
}

func (h *Handler) recreateTovar(c *gin.Context) {
	req := &model.TovarOutput{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	tovar := &model.Tovar{
		TovarID: req.TovarID,
		ShopID:  req.ShopID,
	}

	err = h.services.RecreateTovar(tovar)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) recreateTechCart(c *gin.Context) {
	req := &model.TechCartResponse{}

	err := req.ParseRequest(c)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	techCart := &model.TechCart{
		TechCartID: req.TechCartID,
		ShopID:     req.ShopID,
	}

	err = h.services.RecreateTechCart(techCart)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getToAddTovar(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetToAddTovar(&filter)
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

func (h *Handler) getToAddTechCart(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetToAddTechCart(&filter)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	if filter.Page == 0 {
		sendGeneral(res, c)
	} else {
		sendPagination(filter.Page, pageCount, res, c)
	}
}
