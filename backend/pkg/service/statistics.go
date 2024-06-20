package service

import (
	"time"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type StatisticsService struct {
	repo repository.Statistics
}

func NewStatisticsService(repo repository.Statistics) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetWorkersStat(filter *model.Filter) ([]*model.WorkerStat, int64, error) {

	res, count, err := s.repo.GetWorkersStat(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *StatisticsService) TodayStatistics(filter *model.Filter) (*model.TodayStatistics, error) {

	now := time.Now().UTC()
	nowFormat := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return s.repo.TodayStatistics(nowFormat, filter)
}

func (s *StatisticsService) EveryDayStatistics(filter *model.Filter) (*model.TotalStatistics, error) {
	now := filter.From.Local()
	total := &model.TotalStatistics{}
	nowFormat := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	res, err := s.repo.EveryDayStatistics(nowFormat, filter)
	if err != nil {
		return nil, err
	}
	total.Statistics = res
	for _, stat := range res {
		total.TotalRevenue += stat.Revenue
		total.TotalProfit += stat.Profit
		total.TotalCost += stat.Cost
		total.TotalChecks += stat.Checks
		total.TotalVisitors += stat.Visitors
	}
	if total.TotalChecks != 0 {
		total.TotalAvgCheck = total.TotalRevenue / float32(total.TotalChecks)
	}
	return total, nil
}

func (s *StatisticsService) EveryWeekStatistics(filter *model.Filter) (*model.TotalStatistics, error) {
	now := filter.From.Local()
	total := &model.TotalStatistics{}
	nowFormat := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	res, err := s.repo.EveryWeekStatistics(nowFormat, filter)
	if err != nil {
		return nil, err
	}
	total.Statistics = res
	for _, stat := range res {
		total.TotalRevenue += stat.Revenue
		total.TotalProfit += stat.Profit
		total.TotalCost += stat.Cost
		total.TotalChecks += stat.Checks
		total.TotalVisitors += stat.Visitors
	}
	if total.TotalChecks != 0 {
		total.TotalAvgCheck = total.TotalRevenue / float32(total.TotalChecks)
	}
	return total, nil
}
func (s *StatisticsService) EveryMonthStatistics(filter *model.Filter) (*model.TotalStatistics, error) {
	now := filter.From.Local()
	total := &model.TotalStatistics{}
	nowFormat := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	res, err := s.repo.EveryMonthStatistics(nowFormat, filter)
	if err != nil {
		return nil, err
	}
	total.Statistics = res
	for _, stat := range res {
		total.TotalRevenue += stat.Revenue
		total.TotalProfit += stat.Profit
		total.TotalCost += stat.Cost
		total.TotalChecks += stat.Checks
		total.TotalVisitors += stat.Visitors
	}
	if total.TotalChecks != 0 {
		total.TotalAvgCheck = total.TotalRevenue / float32(total.TotalChecks)
	}
	return total, nil
}

func (s *StatisticsService) Payments(filter *model.Filter) (*model.GlobalPayment, int64, error) {
	res, count, err := s.repo.Payments(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.ChecksPageSize)
	payments := &model.GlobalPayment{
		Payments: res,
	}
	for _, payment := range res {
		payments.TotalTotal += payment.Total
		payments.TotalCard += payment.Card
		payments.TotalCash += payment.Cash
		payments.TotalCheckCount += payment.CheckCount
		payments.TotalCheckCard += payment.CheckCard
		payments.TotalCheckCash += payment.CheckCash
		payments.TotalCheckMixed += payment.CheckMixed
		payments.TotalCheckZero += payment.CheckZero
	}
	return payments, pageCount, nil
}

func (s *StatisticsService) DaysOfTheWeek(filter *model.Filter) (*model.DaysOfTheWeek, error) {
	res, err := s.repo.DaysOfTheWeek(filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *StatisticsService) StatByHour(filter *model.Filter) ([]*model.StatByHour, error) {
	res, err := s.repo.StatByHour(filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *StatisticsService) ABC(filter *model.Filter) ([]*model.ABC, error) {
	res, err := s.repo.ABC(filter)
	if err != nil {
		return nil, err
	}

	type Total struct {
		TotalRevenue float32
		TotalProfit  float32
		TotalSales   int
	}

	total := &Total{}

	for _, val := range res {
		if val.Revenue > 0 {
			total.TotalRevenue += val.Revenue
		}
		if val.Profit > 0 {
			total.TotalProfit += val.Profit
		}
		if val.Sales > 0 {
			total.TotalSales += val.Sales
		}
	}

	for _, val := range res {
		if val.Revenue > 0 {
			val.RevenuePercent = val.Revenue / total.TotalRevenue * 100
		}
		if val.Profit > 0 {
			val.ProfitPercent = val.Profit / total.TotalProfit * 100
		}
		if val.Sales > 0 {
			val.SalesPercent = float32(val.Sales) / float32(total.TotalSales) * 100
		}
	}

	type Percent struct {
		Revenue float32
		Profit  float32
		Sales   float32
	}

	percent := &Percent{}

	for _, val := range res {
		val.ProfitLetter = "C"
		val.RevenueLetter = "C"
		val.SalesLetter = "C"
		if val.Revenue > 0 {
			percent.Revenue += val.RevenuePercent
			if percent.Revenue <= 80 {
				val.RevenueLetter = "A"
			} else if percent.Revenue <= 95 {
				val.RevenueLetter = "B"
			}
		}
		if val.Profit > 0 {
			percent.Profit += val.ProfitPercent
			if percent.Profit <= 80 {
				val.ProfitLetter = "A"
			} else if percent.Profit <= 95 {
				val.ProfitLetter = "B"
			}
		}
		if val.Sales > 0 {
			percent.Sales += val.SalesPercent
			if percent.Sales <= 80 {
				val.SalesLetter = "A"
			} else if percent.Sales <= 95 {
				val.SalesLetter = "B"
			}
		}
	}

	return res, nil
}

func (s *StatisticsService) TopSales(filter *model.Filter) ([]*model.ItemOutput, error) {
	res, err := s.repo.TopSales(filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}
