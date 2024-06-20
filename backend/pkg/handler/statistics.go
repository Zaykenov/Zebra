package handler

import (
	"zebra/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getWorkersStat(c *gin.Context) {

	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}

	res, pageCount, err := h.services.Statistics.GetWorkersStat(&filter)

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

func (h *Handler) todayStatistics(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.TodayStatistics(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) everyDayStatistics(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.EveryDayStatistics(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) everyWeekStatistics(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.EveryWeekStatistics(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) everyMonthStatistics(c *gin.Context) {
	filter := &model.Filter{}
	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.EveryMonthStatistics(filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) payments(c *gin.Context) {
	/*page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 0
	}
	sort := c.Query("sort")
	if err != nil {
		sort = "id DESC"
	}
	now := time.Now()
	layout := "2006-01-02"

	from, err := time.Parse(layout, c.Query("from"))
	if err != nil {
		from = now.AddDate(0, -1, 0)
	}
	to, err := time.Parse(layout, c.Query("to"))
	if err != nil {
		to = now
	}*/
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, pageCount, err := h.services.Statistics.Payments(&filter)
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

func (h *Handler) DaysOfTheWeek(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.DaysOfTheWeek(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) StatByHour(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.StatByHour(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) ABC(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.ABC(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}

func (h *Handler) TopSales(c *gin.Context) {
	filter := model.Filter{}

	if err := filter.ParseRequest(c); err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	res, err := h.services.Statistics.TopSales(&filter)
	if err != nil {
		h.defaultErrorHandler(c, err)
		return
	}
	sendGeneral(res, c)
}
