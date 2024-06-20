package handler

import (
	"io/ioutil"
	"net/http"
	"strings"

	"zebra/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

// serverErrorMessage creator
func serverErrorMessage(StatusCode int, Message string) *model.DefaultResponse {
	response := &model.DefaultResponse{}
	response.StatusCode = StatusCode
	response.Data = Message
	return response
}

// defaultErrorHandler error only with status code
func defaultErrorHandler(c *gin.Context, Err error) {
	fullError := Err.Error()

	parts := strings.Split(fullError, "|")
	mainMessage := strings.TrimSpace(parts[0])

	switch mainMessage {

	case "bad request":
		c.AbortWithStatusJSON(http.StatusBadRequest, serverErrorMessage(400, fullError))
	case "wrong username or password":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "not found":
		c.AbortWithStatusJSON(http.StatusNotFound, serverErrorMessage(404, fullError))
	case "incorrect token":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "username is already taken":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "file system":
		c.AbortWithStatusJSON(http.StatusInternalServerError, serverErrorMessage(8001, fullError))
	case "server error":
		c.AbortWithStatusJSON(http.StatusInternalServerError, serverErrorMessage(500, fullError))
	case "open shift":
		c.AbortWithStatusJSON(666, serverErrorMessage(666, fullError))
	case "close shift":
		c.AbortWithStatusJSON(667, serverErrorMessage(667, fullError))
	case "type new netCost for item":
		c.AbortWithStatusJSON(668, serverErrorMessage(668, fullError))
	case "check is closed":
		c.AbortWithStatusJSON(http.StatusLocked, serverErrorMessage(http.StatusLocked, fullError))
	case "close opened checks":
		c.AbortWithStatusJSON(444, serverErrorMessage(444, fullError))
	case "can not add new inventarization between two inventarizations":
		c.AbortWithStatusJSON(445, serverErrorMessage(445, fullError))
	case "access restricted":
		c.AbortWithStatusJSON(http.StatusForbidden, serverErrorMessage(http.StatusForbidden, fullError))
	case "could not find the shift for given shop":
		c.AbortWithStatusJSON(888, serverErrorMessage(888, fullError))
	default:
		c.AbortWithStatusJSON(http.StatusNotFound, serverErrorMessage(8000, fullError))
	}
}

// defaultErrorHandler error only with status code
func (h *Handler) defaultErrorHandler(c *gin.Context, Err error) {
	fullError := Err.Error()
	body, err := ioutil.ReadAll(c.Request.Body)
	bodyString := ""
	if err == nil {
		bodyString = string(body)
	}
	if bodyString == "" {
		bodyString = c.Request.URL.Path
	} else {
		bodyString = bodyString + "\n" + c.Request.URL.Path
	}
	h.services.SaveError(fullError, bodyString)
	parts := strings.Split(fullError, "|")
	mainMessage := strings.TrimSpace(parts[0])

	switch mainMessage {
	case "bad request":
		c.AbortWithStatusJSON(http.StatusBadRequest, serverErrorMessage(400, fullError))
	case "wrong username or password":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "not found":
		c.AbortWithStatusJSON(http.StatusNotFound, serverErrorMessage(404, fullError))
	case "incorrect token":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "username is already taken":
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverErrorMessage(401, fullError))
	case "file system":
		c.AbortWithStatusJSON(http.StatusInternalServerError, serverErrorMessage(8001, fullError))
	case "server error":
		c.AbortWithStatusJSON(http.StatusInternalServerError, serverErrorMessage(500, fullError))
	case "open shift":
		c.AbortWithStatusJSON(666, serverErrorMessage(666, fullError))
	case "close shift":
		c.AbortWithStatusJSON(667, serverErrorMessage(667, fullError))
	case "check is closed":
		c.AbortWithStatusJSON(http.StatusLocked, serverErrorMessage(http.StatusLocked, fullError))
	case "close opened checks":
		c.AbortWithStatusJSON(444, serverErrorMessage(444, fullError))
	case "can not add new inventarization between two inventarizations":
		c.AbortWithStatusJSON(445, serverErrorMessage(445, fullError))
	case "access restricted":
		c.AbortWithStatusJSON(http.StatusForbidden, serverErrorMessage(http.StatusForbidden, fullError))
	default:
		c.AbortWithStatusJSON(http.StatusNotFound, serverErrorMessage(8000, fullError))
	}
}

// sendGeneral sends general data
func sendGeneral(data interface{}, c *gin.Context) {
	gr := model.SuccessResponse()
	gr.Data = data

	c.JSON(http.StatusOK, gr)
}

// sendSuccess sends response success
func sendSuccess(c *gin.Context) {
	gr := &model.DefaultResponse{StatusCode: 200}
	c.JSON(http.StatusOK, gr)
}

func sendPagination(cPage int, tPage int64, data interface{}, c *gin.Context) {
	gr := model.SuccessResponse()

	pagination := &model.DefaultPage{
		TotalPages:  tPage,
		CurrentPage: cPage,
		Data:        data,
	}

	gr.Data = pagination

	c.JSON(http.StatusOK, gr)
}
