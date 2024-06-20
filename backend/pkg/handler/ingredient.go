package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addIngredient(c *gin.Context) {
	req := &model.ReqIngredient{}

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
	ingredientMaster := &model.IngredientMaster{
		Name:     req.Name,
		Category: req.Category,
		Image:    req.Image,
		Measure:  req.Measure,
	}
	if role == utils.MasterRole {
		ingredientMaster.Status = utils.MenuStatusApproved
	} else {
		ingredientMaster.Status = utils.MenuStatusPending
	}
	newIngredient, err := h.services.Master.AddIngredientMaster(ingredientMaster)
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
	ingredients := []*model.Ingredient{}
	for _, shop := range req.ShopID {
		ingredient := &model.Ingredient{
			ShopID:       shop,
			IngredientID: newIngredient.ID,
			Name:         newIngredient.Name,
			Category:     newIngredient.Category,
			Image:        newIngredient.Image,
			Measure:      newIngredient.Measure,
			IsVisible:    true,
			Cost:         req.Cost,
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

func (h *Handler) getAllIngredient(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Ingredient.GetAllIngredient(&filter)

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

func (h *Handler) getIngredient(c *gin.Context) {
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
	res, err := h.services.Ingredient.GetIngredient(id, &filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateIngredient(c *gin.Context) {
	req := model.ReqIngredient{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	shops, err := getAvailableShops(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateIngredient(&req, shops)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteIngredient(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteIngredient(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) addCategoryIngredient(c *gin.Context) {
	req := &model.CategoryIngredient{}

	err := req.ParseRequest(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.AddCategoryIngredient(req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendSuccess(c)
}

func (h *Handler) getAllCategoryIngredient(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Ingredient.GetAllCategoryIngredient(&filter)
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

func (h *Handler) getCategoryIngredient(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Ingredient.GetCategoryIngredient(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateCategoryIngredient(c *gin.Context) {
	req := model.CategoryIngredient{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateCategoryIngredient(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteCategoryIngredient(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteCategoryIngredient(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) addNabor(c *gin.Context) {
	req := &model.Nabor{}

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
	if role != utils.MasterRole {
		h.defaultErrorHandler(c, errors.New("only master can add nabor"))
		return
	}
	ingredientsNaborsMaster := []*model.IngredientNaborMaster{}
	for _, v := range req.Ingredients {
		ingredientsNaborsMaster = append(ingredientsNaborsMaster, &model.IngredientNaborMaster{
			IngredientID: v.IngredientID,
			Brutto:       v.Brutto,
			Price:        v.Price,
		})
	}
	naborMaster := &model.NaborMaster{
		Name:        req.Name,
		Min:         req.Min,
		Max:         req.Max,
		Ingredients: ingredientsNaborsMaster,
		Replaces:    req.Replaces,
	}
	if role == utils.MasterRole {
		naborMaster.Status = utils.MenuStatusApproved
	} else {
		naborMaster.Status = utils.MenuStatusPending
	}

	newNabor, err := h.services.AddNaborMaster(naborMaster)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	naborID := newNabor.ID

	if role == utils.MasterRole {
		shops, err := getAvailableShops(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		req.Shops = shops
	}

	nabor := &model.Nabor{
		Shops:       req.Shops,
		NaborID:     naborID,
		Name:        req.Name,
		Min:         req.Min,
		Max:         req.Max,
		Replaces:    req.Replaces,
		Ingredients: req.Ingredients,
	}
	res, err := h.services.AddNabor(nabor)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) addNabors(c *gin.Context) {
	req := []*model.Nabor{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	nabors := []*model.Nabor{}
	for _, nabor := range req {
		var id int = nabor.ID
		var naborID int = nabor.ID
		if nabor.ID == 0 {
			ingredientsNaborsMaster := []*model.IngredientNaborMaster{}
			for _, v := range nabor.Ingredients {
				ingredientsNaborsMaster = append(ingredientsNaborsMaster, &model.IngredientNaborMaster{
					IngredientID: v.IngredientID,
					NaborID:      v.ID,
					Brutto:       v.Brutto,
					Price:        v.Price,
				})
			}
			naborMaster := &model.NaborMaster{
				ID:          nabor.ID,
				Name:        nabor.Name,
				Min:         nabor.Min,
				Max:         nabor.Max,
				Ingredients: ingredientsNaborsMaster,
				Replaces:    nabor.Replaces,
			}
			if role == utils.MasterRole {
				naborMaster.Status = utils.MenuStatusApproved
			} else {
				naborMaster.Status = utils.MenuStatusPending
			}

			newNabor, err := h.services.AddOrUpdateNaborMaster(naborMaster)
			if err != nil {
				h.defaultErrorHandler(c, err)
				return
			}
			naborID = newNabor.ID
		}
		shops, err := getAvailableShops(c)
		if err != nil {
			h.defaultErrorHandler(c, err)
			return
		}
		nabor.Shops = shops

		shopNabor := &model.Nabor{
			ID:          id,
			NaborID:     naborID,
			Shops:       shops,
			Name:        nabor.Name,
			Min:         nabor.Min,
			Max:         nabor.Max,
			Replaces:    nabor.Replaces,
			Ingredients: nabor.Ingredients,
		}
		nabors = append(nabors, shopNabor)
	}
	_, err = h.services.AddNabors(nabors)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(nabors, c)
}

func (h *Handler) getAllNabor(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.Ingredient.GetAllNabor(filter)
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

func (h *Handler) getNabor(c *gin.Context) {
	keys := c.Request.URL.Query()["id"]
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Ingredient.GetNabor(id)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendGeneral(res, c)
}

func (h *Handler) updateNabor(c *gin.Context) {
	req := model.Nabor{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.UpdateNabor(&req)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) deleteNabor(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	err = h.services.DeleteNabor(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getTechCartByIngredientID(c *gin.Context) {
	req := model.ReqID{}

	err := req.ParseRequest(c)

	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, err := h.services.GetTechCartByIngredientID(req.ID)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) getDeletedIngredient(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetDeletedIngredient(&filter)
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

func (h *Handler) recreateIngredient(c *gin.Context) {
	req := &model.IngredientOutput{}

	err := req.ParseRequest(c)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	ingredient := &model.Ingredient{
		IngredientID: req.IngredientID,
		ShopID:       req.ShopID,
	}

	err = h.services.RecreateIngredient(ingredient)
	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	sendSuccess(c)
}

func (h *Handler) getToAddIngredient(c *gin.Context) {
	filter := model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.GetToAddIngredient(&filter)
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
