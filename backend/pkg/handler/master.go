package handler

import (
	"encoding/json"
	"strconv"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetShopsMaster(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := h.services.GetShopsMaster()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res := []*model.ShopFromMaster{}
	for _, shop := range shops {
		shopMaster := &model.ShopFromMaster{
			ID:   shop.ID,
			Name: shop.Name,
		}
		ids := []int{shop.ID}
		token, err := h.services.User.GenerateToken(id, utils.MasterRole, ids, shop.ID)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		shopMaster.Token = token
		res = append(res, shopMaster)
	}
	sendGeneral(res, c)
}

func (h *Handler) NormaliseTovars(c *gin.Context) {
	err := h.services.NormaliseTovars()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) NormaliseIngredients(c *gin.Context) {
	err := h.services.NormaliseIngredients()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) NormaliseTechCarts(c *gin.Context) {
	err := h.services.NormaliseTechCarts()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) NormaliseNabors(c *gin.Context) {
	err := h.services.NormaliseNabors()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) getAllTovarMaster(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	tovars, pageCount, err := h.services.GetAllTovarMaster(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if filter.Page != 0 {
		sendPagination(filter.Page, pageCount, tovars, c)
	} else {
		sendGeneral(tovars, c)
	}
}

func (h *Handler) getTovarMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	tovar, err := h.services.GetTovarMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(tovar, c)
}

func (h *Handler) getAllIngredientMaster(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	ingredients, pageCount, err := h.services.GetAllIngredientMaster(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	if filter.Page != 0 {
		sendPagination(filter.Page, pageCount, ingredients, c)
	} else {
		sendGeneral(ingredients, c)
	}
}

func (h *Handler) getIngredientMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	ingredient, err := h.services.GetIngredientMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(ingredient, c)
}

func (h *Handler) getAllTechCartMaster(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	techCart, pageCount, err := h.services.GetAllTechCartMaster(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if filter.Page != 0 {
		sendPagination(filter.Page, pageCount, techCart, c)
	} else {
		sendGeneral(techCart, c)
	}
}

func (h *Handler) getTechCartMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	techCart, err := h.services.GetTechCartMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(techCart, c)
}

func (h *Handler) getAllNaborMaster(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	nabors, pageCount, err := h.services.GetAllNaborMaster(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	if filter.Page != 0 {
		sendPagination(filter.Page, pageCount, nabors, c)
	} else {
		sendGeneral(nabors, c)
	}
}

func (h *Handler) getNaborMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.GetNaborMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) UpdateTovarMaster(c *gin.Context) {
	req := model.ReqTovarMaster{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateTovarMaster(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) DeleteTovarMaster(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteTovarMaster(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) ConfirmTovarMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	tovarMaster, err := h.services.ConfirmTovarMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := h.services.GetShopsMaster()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	existTovars, err := h.services.GetIdsOfShopsWhereTheTovarAlreadyExist(tovarMaster.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	var cost float32 = 0
	shopsWithoutExistTovars := []int{}
	for _, shop := range shops {
		exist := false
		for _, existTovar := range existTovars {
			if shop.ID == existTovar {
				exist = true
				break
			}
		}
		if !exist {
			shopsWithoutExistTovars = append(shopsWithoutExistTovars, shop.ID)
		}
	}
	tovars := []*model.Tovar{}
	for _, shop := range shopsWithoutExistTovars {
		tovar := &model.Tovar{
			ShopID:    shop,
			TovarID:   tovarMaster.ID,
			Name:      tovarMaster.Name,
			Category:  tovarMaster.Category,
			Image:     tovarMaster.Image,
			Tax:       tovarMaster.Tax,
			Measure:   tovarMaster.Measure,
			Discount:  tovarMaster.Discount,
			IsVisible: true,
			Price:     tovarMaster.Price,
			Cost:      cost,
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

func (h *Handler) RejectTovarMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.RejectTovarMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) UpdateTechCartMaster(c *gin.Context) {
	req := &model.ReqTechCartMaster{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateTechCartMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) ConfirmTechCartMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	techCartMaster, err := h.services.ConfirmTechCartMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := h.services.GetShopsMaster()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	existTechCarts, err := h.services.GetIdsOfShopsWhereTheTechCartAlreadyExist(techCartMaster.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shopsWithoutExistTechCarts := []int{}
	for _, shop := range shops {
		exist := false
		for _, existTechCart := range existTechCarts {
			if shop.ID == existTechCart {
				exist = true
				break
			}
		}
		if !exist {
			shopsWithoutExistTechCarts = append(shopsWithoutExistTechCarts, shop.ID)
		}
	}
	techCarts := []*model.TechCart{}
	for _, shop := range shopsWithoutExistTechCarts {
		ingredients := []*model.IngredientTechCart{}
		for _, ingredient := range techCartMaster.Ingredients {
			ingredientTechCart := &model.IngredientTechCart{
				IngredientID: ingredient.IngredientID,
				ShopID:       shop,
				TechCartID:   techCartMaster.ID,
				Brutto:       ingredient.Brutto,
			}
			ingredients = append(ingredients, ingredientTechCart)
		}
		nabors := []*model.NaborTechCart{}
		for _, nabor := range techCartMaster.Nabor {
			naborTechCart := &model.NaborTechCart{
				NaborID:    nabor.NaborID,
				ShopID:     shop,
				TechCartID: techCartMaster.ID,
			}
			nabors = append(nabors, naborTechCart)
		}
		techCart := &model.TechCart{
			ShopID:      shop,
			TechCartID:  techCartMaster.ID,
			Name:        techCartMaster.Name,
			Category:    techCartMaster.Category,
			Image:       techCartMaster.Image,
			Tax:         techCartMaster.Tax,
			Measure:     techCartMaster.Measure,
			Discount:    techCartMaster.Discount,
			IsVisible:   true,
			Price:       techCartMaster.Price,
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

func (h *Handler) RejectTechCartMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.RejectTechCartMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) DeleteTechCartMaster(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteTechCartMaster(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) UpdateIngredientMaster(c *gin.Context) {
	req := &model.ReqIngredientMaster{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.UpdateIngredientMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) ConfirmIngredientMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	ingredientMaster, err := h.services.ConfirmIngredientMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := h.services.GetShopsMaster()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	existIngredients, err := h.services.GetIdsOfShopsWhereTheIngredientAlreadyExist(ingredientMaster.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	var cost float32 = 0
	shopsWithoutExistIngredients := []int{}
	for _, shop := range shops {
		exist := false
		for _, existIngredient := range existIngredients {
			if shop.ID == existIngredient {
				exist = true
				break
			}
		}
		if !exist {
			shopsWithoutExistIngredients = append(shopsWithoutExistIngredients, shop.ID)
		}
	}
	ingredients := []*model.Ingredient{}
	for _, shop := range shopsWithoutExistIngredients {
		ingredient := &model.Ingredient{
			ShopID:       shop,
			IngredientID: ingredientMaster.ID,
			Name:         ingredientMaster.Name,
			Category:     ingredientMaster.Category,
			Image:        ingredientMaster.Image,
			Measure:      ingredientMaster.Measure,
			IsVisible:    true,
			Cost:         cost,
		}
		ingredients = append(ingredients, ingredient)
	}
	err = h.services.AddIngredients(ingredients)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) RejectIngredientMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.RejectIngredientMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) DeleteIngredientMaster(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteIngredientMaster(req.ID)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) createNaborMaster(c *gin.Context) {
	req := []*model.NaborMaster{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.CreateNaborMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) updateTechCartsMaster(c *gin.Context) {
	req := []*model.ReqTechCartMaster{} //nice
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err := h.services.UpdateTechCartsMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) addNaborMaster(c *gin.Context) {
	req := &model.NaborMaster{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	_, err := h.services.AddNaborMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) updateNaborMaster(c *gin.Context) {
	req := &model.NaborMaster{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	_, err := h.services.UpdateNaborMaster(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) ConfirmNaborMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	naborMaster, err := h.services.ConfirmNaborMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := h.services.GetShopsMaster()
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	existNabors, err := h.services.GetIdsOfShopsWhereTheNaborAlreadyExist(naborMaster.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shopsWithoutExistNabors := []int{}
	for _, shop := range shops {
		exist := false
		for _, existNabor := range existNabors {
			if shop.ID == existNabor {
				exist = true
				break
			}
		}
		if !exist {
			shopsWithoutExistNabors = append(shopsWithoutExistNabors, shop.ID)
		}
	}
	nabors := []*model.Nabor{}
	for _, shop := range shopsWithoutExistNabors {
		ingredients := []*model.IngredientNabor{}
		for _, ingredient := range naborMaster.Ingredients {
			ingredientNabor := &model.IngredientNabor{
				IngredientID: ingredient.IngredientID,
				ShopID:       shop,
				NaborID:      naborMaster.ID,
				Brutto:       ingredient.Brutto,
				Price:        ingredient.Price,
			}
			ingredients = append(ingredients, ingredientNabor)
		}
		nabor := &model.Nabor{
			ShopID:      shop,
			NaborID:     naborMaster.ID,
			Name:        naborMaster.Name,
			Min:         naborMaster.Min,
			Max:         naborMaster.Max,
			Replaces:    naborMaster.Replaces,
			Ingredients: ingredients,
		}
		nabors = append(nabors, nabor)
	}
	_, err = h.services.AddNabors(nabors)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) RejectNaborMaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	err = h.services.RejectNaborMaster(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) DeleteNaborMaster(c *gin.Context) {
	req := model.ReqID{}
	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteNaborMaster(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}
