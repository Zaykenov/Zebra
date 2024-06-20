package handler

import (
	"errors"
	"net/http"
	"strings"
	"zebra/model"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "role"
	accessShopsCtx      = "accessShops"
	bindShopCtx         = "bindShop"
	IdempotencyKey      = "Idempotency-Key"
	version             = "uuid_version"
)

func (h *Handler) CheckIdempotency(ReqCheck *model.ReqCheck) (int, error) {
	if ReqCheck.IdempotencyKey == "" {
		return 0, nil
	}
	check, err := h.services.GetCheckByIdempotency(ReqCheck.IdempotencyKey)
	if err != nil {
		if err.Error() == "record not found" {
			return 0, nil
		} else {
			return 0, err
		}
	}
	if check == nil {
		return 0, nil
	}
	if ReqCheck.Version <= check.Version {
		return 0, errors.New("idempotency key already exists")
	} else {
		return check.ID, nil
	}

}

func (h *Handler) workerIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, role, accessShops, bindShop, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if role != utils.ManagerRole && role != utils.WorkerRole && role != utils.MasterRole {
		newErrorResponse(c, http.StatusUnauthorized, "access restricted")
		return
	}

	c.Set(userCtx, userId)
	c.Set(roleCtx, role)
	c.Set(accessShopsCtx, accessShops)
	c.Set(bindShopCtx, bindShop)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, role, accessShops, bindShop, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if role != utils.ManagerRole && role != utils.MasterRole {
		newErrorResponse(c, http.StatusUnauthorized, "access restricted")
		return
	}
	c.Set(userCtx, userId)
	c.Set(roleCtx, role)
	c.Set(accessShopsCtx, accessShops)
	c.Set(bindShopCtx, bindShop)
}

func (h *Handler) masterIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, role, accessShops, bindShop, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if role != utils.MasterRole {
		newErrorResponse(c, http.StatusUnauthorized, "access restricted")
		return
	}
	c.Set(userCtx, userId)
	c.Set(roleCtx, role)
	c.Set(accessShopsCtx, accessShops)
	c.Set(bindShopCtx, bindShop)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get(roleCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	idInt, ok := role.(string)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getAvailableShops(c *gin.Context) ([]int, error) {
	shops, ok := c.Get(accessShopsCtx)
	if !ok {
		return nil, errors.New("user id not found")
	}

	idInt, ok := shops.([]int)
	if !ok {
		return nil, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getBindShop(c *gin.Context) (int, error) {
	bindShopCtx, ok := c.Get(bindShopCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := bindShopCtx.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
