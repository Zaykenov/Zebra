package repository

import (
	"database/sql"
	"fmt"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type StatisticsDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewStatisticsDB(db *sql.DB, gormDB *gorm.DB) *StatisticsDB {
	return &StatisticsDB{db: db, gormDB: gormDB}
}

func (r *StatisticsDB) GetWorkersStat(filter *model.Filter) ([]*model.WorkerStat, int64, error) {
	workerStat := []*model.WorkerStat{}

	res := r.gormDB.Model(model.Worker{}).Select("workers.id, workers.name, COUNT(checks.id) AS check_num, SUM(checks.cost) as cost, SUM(checks.sum) as revenue").Joins("LEFT JOIN checks ON workers.id = checks.worker_id").Where("checks.status = ? and checks.shop_id IN (?)", utils.StatusClosed, filter.AccessibleShops).Group("workers.id")

	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.Worker{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}
	if err := newRes.Scan(&workerStat).Error; err != nil {
		return nil, 0, err
	}

	for _, worker := range workerStat {

		worker.Profit = worker.Revenue - worker.Cost
		if worker.CheckNum != 0 {
			worker.AvgCheck = worker.Revenue / float32(worker.CheckNum)
		}
	}
	return workerStat, count, nil
}

func (r *StatisticsDB) TodayStatistics(timeStamp time.Time, filter *model.Filter) (*model.TodayStatistics, error) {
	statistics := &model.Statistics{}

	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors FROM checks WHERE DATE(checks.closed_at) = ?::date and checks.status = ? and checks.shop_id IN (?)", timeStamp.Format("2006-01-02"), utils.StatusClosed, filter.Shop).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors FROM checks WHERE DATE(checks.closed_at) = ?::date and checks.status = ? and checks.shop_id IN (?)", timeStamp.Format("2006-01-02"), utils.StatusClosed, filter.AccessibleShops).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	}

	statistics.Profit = statistics.Revenue - statistics.Cost
	if statistics.Checks != 0 {
		statistics.AvgCheck = statistics.Revenue / float32(statistics.Checks)
	}
	statistics.From = timeStamp
	statistics.To = timeStamp

	prevStatistics, err := r.PrevDayHour(filter)
	if err != nil {
		return nil, err
	}

	statisticsResponse := &model.TodayStatistics{}

	statisticsResponse.Revenue = statistics.Revenue
	statisticsResponse.PrevRevenue = prevStatistics.Revenue

	if statisticsResponse.Revenue != 0 && statisticsResponse.PrevRevenue != 0 {
		statisticsResponse.PercentRevenue = int((statisticsResponse.Revenue*100)/statisticsResponse.PrevRevenue) - 100
	}

	statisticsResponse.Profit = statistics.Profit
	statisticsResponse.PrevProfit = prevStatistics.Profit

	if statisticsResponse.Profit != 0 && statisticsResponse.PrevProfit != 0 {
		statisticsResponse.PercentProfit = int((statisticsResponse.Profit*100)/statisticsResponse.PrevProfit) - 100
	}

	statisticsResponse.Checks = statistics.Checks
	statisticsResponse.PrevChecks = prevStatistics.Checks

	if statisticsResponse.Checks != 0 && statisticsResponse.PrevChecks != 0 {
		statisticsResponse.PercentChecks = int((float64(statisticsResponse.Checks)*100)/float64(statisticsResponse.PrevChecks)) - 100
	}

	statisticsResponse.Visitors = statistics.Visitors
	statisticsResponse.PrevVisitors = prevStatistics.Visitors

	if statisticsResponse.Visitors != 0 && statisticsResponse.PrevVisitors != 0 {
		statisticsResponse.PercentVisitors = int((float64(statisticsResponse.Visitors)*100)/float64(statisticsResponse.PrevVisitors)) - 100
	}

	statisticsResponse.AvgCheck = statistics.AvgCheck
	statisticsResponse.PrevAvgCheck = prevStatistics.AvgCheck

	if statisticsResponse.AvgCheck != 0 && statisticsResponse.PrevAvgCheck != 0 {
		statisticsResponse.PercentAvgCheck = int((statisticsResponse.AvgCheck*100)/statisticsResponse.PrevAvgCheck) - 100
	}

	return statisticsResponse, nil
}

//SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, DATE(checks.closed_at) FROM checks GROUP BY DATE(checks.closed_at)

func (r *StatisticsDB) EveryDayStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error) {
	statistics := []*model.Statistics{}
	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, DATE(checks.closed_at) as from, DATE(checks.closed_at) as to FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?) GROUP BY DATE(checks.closed_at)", time.Format("2006-01-02"), utils.StatusClosed, filter.Shop).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, DATE(checks.closed_at) as from, DATE(checks.closed_at) as to FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?) GROUP BY DATE(checks.closed_at)", time.Format("2006-01-02"), utils.StatusClosed, filter.AccessibleShops).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	}

	for _, stat := range statistics {
		stat.Profit = stat.Revenue - stat.Cost
		if stat.Checks != 0 {
			stat.AvgCheck = stat.Revenue / float32(stat.Checks)
		}
	}
	month := filter.To.Local()

	resStatistics := []*model.Statistics{}

	for time.Before(month) || time.Equal(month) {
		flag := false
		for _, stat := range statistics {
			if stat.From.Day() == time.Day() && stat.From.Month() == time.Month() && stat.From.Year() == time.Year() {
				resStatistics = append(resStatistics, stat)
				flag = true
				break
			}
		}
		if !flag {
			resStatistics = append(resStatistics, &model.Statistics{
				From:     time,
				To:       time,
				Checks:   0,
				Cost:     0,
				Revenue:  0,
				Visitors: 0,
				Profit:   0,
				AvgCheck: 0,
			})
		}
		time = time.AddDate(0, 0, 1)
	}

	return resStatistics, nil
}

func (r *StatisticsDB) EveryWeekStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error) {
	statistics := []*model.Statistics{}

	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors,DATE_PART('week',checks.closed_at::date) as week FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?) GROUP BY DATE_PART('week',checks.closed_at::date)", time.Format("2006-01-02"), utils.StatusClosed, filter.Shop).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors,DATE_PART('week',checks.closed_at::date) as week FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?) GROUP BY DATE_PART('week',checks.closed_at::date)", time.Format("2006-01-02"), utils.StatusClosed, filter.AccessibleShops).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	}
	for _, stat := range statistics {
		stat.Profit = stat.Revenue - stat.Cost
		if stat.Checks != 0 {
			stat.AvgCheck = stat.Revenue / float32(stat.Checks)
		}
	}

	month := filter.To.Local()
	_, startWeek := time.ISOWeek()
	_, endWeek := month.ISOWeek()
	resStatistics := []*model.Statistics{}

	for startWeek <= endWeek {
		flag := false
		for _, stat := range statistics {
			if stat.Week == startWeek {
				resStatistics = append(resStatistics, stat)
				flag = true
				break
			}
		}
		if !flag {
			resStatistics = append(resStatistics, &model.Statistics{
				From:     time,
				To:       time,
				Checks:   0,
				Cost:     0,
				Revenue:  0,
				Visitors: 0,
				Profit:   0,
				AvgCheck: 0,
				Week:     startWeek,
			})
		}
		startWeek++
	}
	i := 0

	for time.Before(month) || time.Equal(month) {
		_, week := time.ISOWeek()
		if resStatistics[i].Week == week {
			resStatistics[i].From = time
			for time.Weekday() != 0 && time.Before(month) {
				time = time.AddDate(0, 0, 1)
			}
			resStatistics[i].To = time
			i++
			if i >= len(resStatistics) {
				i = len(resStatistics) - 1
			}
		}
		time = time.AddDate(0, 0, 1)
	}
	return resStatistics, nil
}

func (r *StatisticsDB) EveryMonthStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error) {
	statistics := []*model.Statistics{}
	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, DATE_PART('month',checks.closed_at::date) as month FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?)  GROUP BY DATE_PART('month',checks.closed_at::date)", time.Format("2006-01-02"), utils.StatusClosed, filter.Shop).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, DATE_PART('month',checks.closed_at::date) as month FROM checks WHERE DATE(checks.closed_at) >= ?::date and checks.status = ? and checks.shop_id IN (?)  GROUP BY DATE_PART('month',checks.closed_at::date)", time.Format("2006-01-02"), utils.StatusClosed, filter.AccessibleShops).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	}
	for _, stat := range statistics {
		stat.Profit = stat.Revenue - stat.Cost
		if stat.Checks != 0 {
			stat.AvgCheck = stat.Revenue / float32(stat.Checks)
		}
	}
	month := filter.To.Local()
	newTime := time
	resStatistics := []*model.Statistics{}
	for newTime.Before(month) || newTime.Equal(month) {
		flag := false
		for _, stat := range statistics {
			if stat.Month == int(newTime.Month()) {
				resStatistics = append(resStatistics, stat)
				flag = true
				break
			}
		}
		if !flag {
			resStatistics = append(resStatistics, &model.Statistics{
				From:     time,
				To:       time,
				Checks:   0,
				Cost:     0,
				Revenue:  0,
				Visitors: 0,
				Profit:   0,
				AvgCheck: 0,
				Month:    int(newTime.Month()),
			})
		}
		newTime = newTime.AddDate(0, 1, 0)
	}
	i := 0
	for time.Before(month) || time.Equal(month) {
		if resStatistics[i].Month == int(time.Month()) {
			resStatistics[i].From = time
			for time.AddDate(0, 0, 1).Day() != 1 && time.Before(month) {
				time = time.AddDate(0, 0, 1)
			}
			resStatistics[i].To = time
			i++
			if i >= len(resStatistics) {
				i = len(resStatistics) - 1
			}
		}
		time = time.AddDate(0, 0, 1)
	}
	return resStatistics, nil
}

func (r *StatisticsDB) Payments(filter *model.Filter) ([]*model.Payment, int64, error) {
	/*payments := []*model.Payment{}
	statistics := []*model.PaymentRead{}
	to := filter.To
	from := filter.From
	res := r.gormDB.Model(model.Check{}).Select("COUNT(checks.id) AS check_count,SUM(checks.sum-checks.discount) AS revenue,  DATE(checks.closed_at) as time, checks.payment").Where("checks.status = ? and DATE(checks.closed_at) <= ?::date and DATE(checks.closed_at) >= ?::date", utils.StatusClosed, to, from).Group("DATE(checks.closed_at), checks.payment").Scan(&statistics)
	if err := res.Error; err != nil {
		return nil, err
	}
	for to.After(from) || to.Equal(from) {
		payments = append(payments, &model.Payment{
			Time:       to,
			CheckCount: 0,
			Card:       0,
			Cash:       0,
			Total:      0,
		})
		to = to.AddDate(0, 0, -1)
	}
	for _, payment := range payments {
		for _, stat := range statistics {
			if stat.Time.Day() == payment.Time.Day() && stat.Time.Month() == payment.Time.Month() && stat.Time.Year() == payment.Time.Year() {
				payment.CheckCount += stat.CheckCount
				payment.Total += stat.Revenue
				if stat.Payment == utils.PaymentCard {
					payment.Card = stat.Revenue
				} else {
					payment.Cash = stat.Revenue
				}
			}
		}S
	}
	return payments, nil*/
	payments := []*model.Payment{}
	check := model.Check{}
	res := r.gormDB.Table("checks").Select("COUNT(checks.id) AS check_count, SUM(checks.cash) AS cash, SUM(checks.card) AS card, SUM(checks.sum-checks.discount) as total, DATE(checks.closed_at) as time, SUM(CASE WHEN checks.cash = 0 AND checks.card != 0 THEN 1 ELSE 0 END) as check_card, SUM(CASE WHEN checks.card = 0 AND checks.cash != 0 THEN 1 ELSE 0 END) as check_cash, SUM(CASE WHEN checks.cash != 0 AND checks.card != 0 THEN 1 ELSE 0 END) as check_mixed, SUM(CASE WHEN checks.cash = 0 AND checks.card = 0 THEN 1 ELSE 0 END) as check_zero").Where("checks.status = ?", utils.StatusClosed).Group("DATE(checks.closed_at)")
	newRes, count, err := filter.FilterResults(res, check, utils.ChecksPageSize, "closed_at", "", "time desc")
	if err != nil {
		return nil, 0, err
	}
	if newRes.Scan(&payments).Error != nil {
		return nil, 0, newRes.Error
	}
	return payments, count, nil
}

//COUNT(checks.id) AS check_count, SUM(checks.cash-checks.discount) AS cash, SUM(checks.card-checks.discount) AS card,  DATE(checks.closed_at) as time

func (r *StatisticsDB) DaysOfTheWeek(filter *model.Filter) (*model.DaysOfTheWeek, error) {
	weeklyDay := &model.DaysOfTheWeek{}
	statistics := []*model.Statistics{}
	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, checks.closed_at::date as from FROM checks WHERE DATE(checks.closed_at) >= ?::date and DATE(checks.closed_at) <= ?::date and checks.status = ? and checks.shop_id IN (?) group by checks.closed_at::date", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02"), utils.StatusClosed, filter.Shop).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT COUNT(checks.id) AS checks, SUM(checks.cost) AS cost, SUM(checks.sum - checks.discount) AS revenue, COUNT(DISTINCT checks.user_id) AS visitors, checks.closed_at::date as from FROM checks WHERE DATE(checks.closed_at) >= ?::date and DATE(checks.closed_at) <= ?::date and checks.status = ? and checks.shop_id IN (?) group by checks.closed_at::date", filter.From.Format("2006-01-02"), filter.To.Format("2006-01-02"), utils.StatusClosed, filter.AccessibleShops).Scan(&statistics).Error; err != nil {
			return nil, err
		}
	}

	var w1 float32 = 0
	var w2 float32 = 0
	var w3 float32 = 0
	var w4 float32 = 0
	var w5 float32 = 0
	var w6 float32 = 0
	var w7 float32 = 0
	for _, stat := range statistics {
		stat.Profit = stat.Revenue - stat.Cost
		if stat.Checks != 0 {
			stat.AvgCheck = stat.Revenue / float32(stat.Checks)
		}
		dayOfTheWeek := stat.From.Weekday()
		switch dayOfTheWeek {
		case time.Monday:
			w1++
			weeklyDay.Monday.Revenue += stat.Revenue
			weeklyDay.Monday.Profit += stat.Profit
			weeklyDay.Monday.Visitors += stat.Visitors
			weeklyDay.Monday.Checks += stat.Checks
			weeklyDay.Monday.AvgCheck += stat.AvgCheck
		case time.Tuesday:
			w2++
			weeklyDay.Tuesday.Revenue += stat.Revenue
			weeklyDay.Tuesday.Profit += stat.Profit
			weeklyDay.Tuesday.Visitors += stat.Visitors
			weeklyDay.Tuesday.Checks += stat.Checks
			weeklyDay.Tuesday.AvgCheck += stat.AvgCheck
		case time.Wednesday:
			w3++
			weeklyDay.Wednesday.Revenue += stat.Revenue
			weeklyDay.Wednesday.Profit += stat.Profit
			weeklyDay.Wednesday.Visitors += stat.Visitors
			weeklyDay.Wednesday.Checks += stat.Checks
			weeklyDay.Wednesday.AvgCheck += stat.AvgCheck
		case time.Thursday:
			w4++
			weeklyDay.Thursday.Revenue += stat.Revenue
			weeklyDay.Thursday.Profit += stat.Profit
			weeklyDay.Thursday.Visitors += stat.Visitors
			weeklyDay.Thursday.Checks += stat.Checks
			weeklyDay.Thursday.AvgCheck += stat.AvgCheck
		case time.Friday:
			w5++
			weeklyDay.Friday.Revenue += stat.Revenue
			weeklyDay.Friday.Profit += stat.Profit
			weeklyDay.Friday.Visitors += stat.Visitors
			weeklyDay.Friday.Checks += stat.Checks
			weeklyDay.Friday.AvgCheck += stat.AvgCheck
		case time.Saturday:
			w6++
			weeklyDay.Saturday.Revenue += stat.Revenue
			weeklyDay.Saturday.Profit += stat.Profit
			weeklyDay.Saturday.Visitors += stat.Visitors
			weeklyDay.Saturday.Checks += stat.Checks
			weeklyDay.Saturday.AvgCheck += stat.AvgCheck
		case time.Sunday:
			w7++
			weeklyDay.Sunday.Revenue += stat.Revenue
			weeklyDay.Sunday.Profit += stat.Profit
			weeklyDay.Sunday.Visitors += stat.Visitors
			weeklyDay.Sunday.Checks += stat.Checks
			weeklyDay.Sunday.AvgCheck += stat.AvgCheck
		}
	}
	for i := 1; i <= 7; i++ {
		switch i {
		case 1:
			if w1 != 0 {
				weeklyDay.Monday.AvgCheck /= w1
			}
		case 2:
			if w2 != 0 {
				weeklyDay.Tuesday.AvgCheck /= w2
			}
		case 3:
			if w3 != 0 {
				weeklyDay.Wednesday.AvgCheck /= w3
			}
		case 4:
			if w4 != 0 {
				weeklyDay.Thursday.AvgCheck /= w4
			}
		case 5:
			if w5 != 0 {
				weeklyDay.Friday.AvgCheck /= w5
			}
		case 6:
			if w6 != 0 {
				weeklyDay.Saturday.AvgCheck /= w6
			}
		case 7:
			if w7 != 0 {
				weeklyDay.Sunday.AvgCheck /= w7
			}
		}
	}
	return weeklyDay, nil
}

func (r *StatisticsDB) StatByHour(filter *model.Filter) ([]*model.StatByHour, error) {
	hour := []*model.StatByHour{}
	now := time.Now().UTC()
	from := time.Date(now.Year(), now.Month(), now.Day()-1, 18, 0, 0, 0, time.UTC)
	to := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT SUM(checks.sum) AS sum, SUM(checks.sum - checks.discount) AS revenue, SUM(checks.cost) AS cost, COUNT(DISTINCT checks.user_id) AS visitors, COUNT(checks.id) AS checks, DATE_TRUNC('hour', checks.closed_at) AS hour FROM checks WHERE checks.closed_at >= ? and checks.closed_at <= ? and checks.status = ? and checks.shop_id IN (?) group by hour", from, to, utils.StatusClosed, filter.Shop).Scan(&hour).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT SUM(checks.sum) AS sum, SUM(checks.sum - checks.discount) AS revenue, SUM(checks.cost) AS cost, COUNT(DISTINCT checks.user_id) AS visitors, COUNT(checks.id) AS checks, DATE_TRUNC('hour', checks.closed_at) AS hour FROM checks WHERE checks.closed_at >= ? and checks.closed_at <= ? and checks.status = ? and checks.shop_id IN (?) group by hour", from, to, utils.StatusClosed, filter.AccessibleShops).Scan(&hour).Error; err != nil {
			return nil, err
		}
	}
	for _, h := range hour {
		h.Profit = h.Revenue - h.Cost
		if h.Checks != 0 {
			h.AvgCheck = h.Revenue / float32(h.Checks)
		}
		h.Hour = h.Hour.UTC()
	}
	hourOut := []*model.StatByHour{}
	length := 0
	for i := 0; i < 24; i++ {
		if i < 6 {
			if length < len(hour) && hour[length].Hour.Hour() == 18+i {
				hourOut = append(hourOut, hour[length])
				length++
			} else {
				hourOut = append(hourOut, &model.StatByHour{
					Hour: time.Date(now.Year(), now.Month(), now.Day()-1, 18+i, 0, 0, 0, time.UTC),
				})
			}
		} else {
			if length < len(hour) && hour[length].Hour.Hour() == i-6 {
				hourOut = append(hourOut, hour[length])
				length++
			} else {
				hourOut = append(hourOut, &model.StatByHour{
					Hour: time.Date(now.Year(), now.Month(), now.Day(), i-6, 0, 0, 0, time.UTC),
				})
			}
		}

	}
	return hourOut, nil
}

func (r *StatisticsDB) PrevDayHour(filter *model.Filter) (*model.TodayStatistics, error) {
	hour := []*model.StatByHour{}
	prevDayStat := &model.TodayStatistics{}

	var from, to time.Time
	now := time.Now().UTC()
	if now.Hour() >= 18 {
		from = time.Date(now.Year(), now.Month(), now.Day()-1, 18, 0, 0, 0, time.UTC)
		to = time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
	} else {
		from = time.Date(now.Year(), now.Month(), now.Day()-2, 18, 0, 0, 0, time.UTC)
		to = time.Date(now.Year(), now.Month(), now.Day()-1, 18, 0, 0, 0, time.UTC)
	}
	if len(filter.Shop) != 0 {
		if err := r.gormDB.Raw("SELECT SUM(checks.sum) AS sum, SUM(checks.sum - checks.discount) AS revenue, SUM(checks.cost) AS cost, COUNT(DISTINCT checks.user_id) AS visitors, COUNT(checks.id) AS checks, DATE_TRUNC('hour', checks.closed_at) AS hour FROM checks WHERE checks.closed_at >= ? and checks.closed_at <= ? and checks.status = ? and checks.shop_id IN (?) group by hour", from, to, utils.StatusClosed, filter.Shop).Scan(&hour).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.gormDB.Raw("SELECT SUM(checks.sum) AS sum, SUM(checks.sum - checks.discount) AS revenue, SUM(checks.cost) AS cost, COUNT(DISTINCT checks.user_id) AS visitors, COUNT(checks.id) AS checks, DATE_TRUNC('hour', checks.closed_at) AS hour FROM checks WHERE checks.closed_at >= ? and checks.closed_at <= ? and checks.status = ? and checks.shop_id IN (?) group by hour", from, to, utils.StatusClosed, filter.AccessibleShops).Scan(&hour).Error; err != nil {
			return nil, err
		}
	}

	for _, stat := range hour {
		timeStamp := stat.Hour.AddDate(0, 0, 1)
		if timeStamp.After(now) {
			break
		}
		stat.Profit = stat.Revenue - stat.Cost
		if stat.Checks != 0 {
			stat.AvgCheck = stat.Revenue / float32(stat.Checks)
		}
		prevDayStat.Revenue += stat.Revenue
		prevDayStat.Profit += stat.Profit
		prevDayStat.Visitors += stat.Visitors
		prevDayStat.Checks += stat.Checks
		prevDayStat.AvgCheck = stat.AvgCheck
	}
	if prevDayStat.Checks != 0 {
		prevDayStat.AvgCheck = prevDayStat.Revenue / float32(prevDayStat.Checks)
	}
	return prevDayStat, nil
}

func (r *StatisticsDB) ABC(filter *model.Filter) ([]*model.ABC, error) {
	abc := []*model.ABC{}
	res := r.gormDB
	if len(filter.Shop) != 0 {
		if filter.Category != 0 {
			res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, MAX(daily_statistics.item_name) AS item_name, daily_statistics.type AS item_type, MAX(CASE WHEN daily_statistics.type = 'techCart' THEN tech_carts.category ELSE tovars.category END) AS item_category, SUM(daily_statistics.sales) AS sales, SUM(daily_statistics.check_price) AS revenue, SUM(daily_statistics.check_price - daily_statistics.check_cost) AS profit").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.Shop, filter.Shop).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ? and (tovars.category = ? or tech_carts.category = ?)", utils.TypeIngredient, filter.From, filter.To, filter.Category, filter.Category).Group("daily_statistics.item_id, daily_statistics.type")
			if res.Error != nil {
				return nil, res.Error
			}
		} else {
			res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, MAX(daily_statistics.item_name) AS item_name, daily_statistics.type AS item_type, MAX(CASE WHEN daily_statistics.type = 'techCart' THEN tech_carts.category ELSE tovars.category END) AS item_category, SUM(daily_statistics.sales) AS sales, SUM(daily_statistics.check_price) AS revenue, SUM(daily_statistics.check_price - daily_statistics.check_cost) AS profit").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.Shop, filter.Shop).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ?", utils.TypeIngredient, filter.From, filter.To).Group("daily_statistics.item_id, daily_statistics.type")
			if res.Error != nil {
				return nil, res.Error
			}
		}

	} else {
		if filter.Category != 0 {
			res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, MAX(daily_statistics.item_name) AS item_name, daily_statistics.type AS item_type, MAX(CASE WHEN daily_statistics.type = 'techCart' THEN tech_carts.category ELSE tovars.category END) AS item_category, SUM(daily_statistics.sales) AS sales, SUM(daily_statistics.check_price) AS revenue, SUM(daily_statistics.check_price - daily_statistics.check_cost) AS profit").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.AccessibleShops, filter.AccessibleShops).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ? and (tovars.category = ? or tech_carts.category = ?)", utils.TypeIngredient, filter.From, filter.To, filter.Category, filter.Category).Group("daily_statistics.item_id, daily_statistics.type")
			if res.Error != nil {
				return nil, res.Error
			}
		} else {
			res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, MAX(daily_statistics.item_name) AS item_name, daily_statistics.type AS item_type, MAX(CASE WHEN daily_statistics.type = 'techCart' THEN tech_carts.category ELSE tovars.category END) AS item_category, SUM(daily_statistics.sales) AS sales, SUM(daily_statistics.check_price) AS revenue, SUM(daily_statistics.check_price - daily_statistics.check_cost) AS profit").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.AccessibleShops, filter.AccessibleShops).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ?", utils.TypeIngredient, filter.From, filter.To).Group("daily_statistics.item_id, daily_statistics.type")
			if res.Error != nil {
				return nil, res.Error
			}
		}
	}

	newRes, _, err := filter.FilterResults(res, model.DailyStatistic{}, utils.NoNeedPageSize, "", fmt.Sprintf("daily_statistics.item_name ilike '%%%s%%'", filter.Search), "SUM(daily_statistics.check_price) desc")
	if err != nil {
		return nil, err
	}

	if err := newRes.Debug().Scan(&abc).Error; err != nil {
		return nil, err
	}

	return abc, nil
}

func (r *StatisticsDB) TopSales(filter *model.Filter) ([]*model.ItemOutput, error) {
	tovars := []*model.ItemOutput{}
	res := r.gormDB
	if len(filter.Shop) != 0 {
		res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id as id, MAX(daily_statistics.item_name) AS name, daily_statistics.type AS type, SUM(daily_statistics.sales) AS sales").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.Shop, filter.Shop).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ?", utils.TypeIngredient, filter.From, filter.To).Group("daily_statistics.item_id, daily_statistics.type").Limit(8)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id as id, MAX(daily_statistics.item_name) AS name, daily_statistics.type AS type, SUM(daily_statistics.sales) AS sales").Joins("left join tech_carts on (daily_statistics.item_id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?) and daily_statistics.type = 'techCart') left join tovars on (daily_statistics.item_id = tovars.tovar_id and tovars.shop_id IN (?) and daily_statistics.type = 'tovar')", filter.AccessibleShops, filter.AccessibleShops).Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ?", utils.TypeIngredient, filter.From, filter.To).Group("daily_statistics.item_id, daily_statistics.type").Limit(8)
		if res.Error != nil {
			return nil, res.Error
		}
	}

	newRes, _, err := filter.FilterResults(res, model.DailyStatistic{}, utils.NoNeedPageSize, "", fmt.Sprintf("daily_statistics.item_name ilike '%%%s%%'", filter.Search), "SUM(daily_statistics.sales) desc")
	if err != nil {
		return nil, err
	}

	if err := newRes.Debug().Scan(&tovars).Error; err != nil {
		return nil, err
	}

	return tovars, nil
}
