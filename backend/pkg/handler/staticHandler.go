package handler

import (
	"errors"
	"os"

	"zebra/pkg/service"
	"zebra/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StaticHandler struct {
	services *service.Service
}

func NewStaticHandler(services *service.Service) *StaticHandler {
	return &StaticHandler{services: services}
}

func (h *StaticHandler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/upload", h.uploadImage)
	router.GET("/itemImage/:fileName", h.ItemImageHandler)
	return router
}

func (h *StaticHandler) ItemImageHandler(c *gin.Context) {
	filename := c.Param("fileName")

	if filename == "" {
		defaultErrorHandler(c, errors.New("bad request"))
		return
	}

	locationPrefix := os.Getenv("LocationItemDocker")

	logrus.Print(locationPrefix + filename)

	c.File(locationPrefix + filename)
}

func (h *StaticHandler) uploadImage(c *gin.Context) {
	imageName, err := utils.CreateItemImage(c)

	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	sendGeneral(imageName, c)
}

func (h *Handler) ItemImageHandler(c *gin.Context) {
	filename := c.Param("fileName")

	if filename == "" {
		defaultErrorHandler(c, errors.New("bad request"))
		return
	}

	locationPrefix := os.Getenv("LocationItemDocker")

	logrus.Print(locationPrefix + filename)

	c.File(locationPrefix + filename)
}

func (h *Handler) uploadImage(c *gin.Context) {
	imageName, err := utils.CreateItemImage(c)

	if err != nil {
		defaultErrorHandler(c, err)
		return
	}

	sendGeneral(imageName, c)
}
