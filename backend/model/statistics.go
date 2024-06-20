package model

import "time"

type Statistics struct {
	Revenue  float32   `json:"revenue"`
	Profit   float32   `json:"profit"`
	Cost     float32   `json:"cost"`
	Checks   int       `json:"checks"`
	Visitors int       `json:"visitors"`
	AvgCheck float32   `json:"avg_check"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
	Week     int       `json:"week"`
	Month    int       `json:"month"`
}

type TotalStatistics struct {
	Statistics    []*Statistics `json:"statistics"`
	TotalRevenue  float32       `json:"total_revenue"`
	TotalProfit   float32       `json:"total_profit"`
	TotalCost     float32       `json:"total_cost"`
	TotalChecks   int           `json:"total_checks"`
	TotalVisitors int           `json:"total_visitors"`
	TotalAvgCheck float32       `json:"total_avg_check"`
}

type GlobalPayment struct {
	TotalCheckCount int        `json:"total_check_count"`
	TotalCheckCard  int        `json:"total_check_card"`
	TotalCheckCash  int        `json:"total_check_cash"`
	TotalCheckMixed int        `json:"total_check_mixed"`
	TotalCheckZero  int        `json:"total_check_zero"`
	TotalCard       float32    `json:"total_card"`
	TotalCash       float32    `json:"total_cash"`
	TotalTotal      float32    `json:"total_total"`
	Payments        []*Payment `json:"payments"`
}

type Payment struct {
	Time       time.Time `json:"time"`
	CheckCount int       `json:"check_count"`
	CheckCard  int       `json:"check_card"`
	CheckCash  int       `json:"check_cash"`
	CheckMixed int       `json:"check_mixed"`
	CheckZero  int       `json:"check_zero"`
	Card       float32   `json:"card"`
	Cash       float32   `json:"cash"`
	Total      float32   `json:"total"`
}

type PaymentRead struct {
	Time       time.Time `json:"time"`
	CheckCount int       `json:"check_count"`
	Revenue    float32   `json:"revenue"`
	Payment    string    `json:"payment"`
}

type TodayStatistics struct {
	Revenue         float32 `json:"revenue"`
	PrevRevenue     float32 `json:"prev_revenue"`
	PercentRevenue  int     `json:"percent_revenue"`
	Profit          float32 `json:"profit"`
	PrevProfit      float32 `json:"prev_profit"`
	PercentProfit   int     `json:"percent_profit"`
	Checks          int     `json:"checks"`
	PrevChecks      int     `json:"prev_checks"`
	PercentChecks   int     `json:"percent_checks"`
	Visitors        int     `json:"visitors"`
	PrevVisitors    int     `json:"prev_visitors"`
	PercentVisitors int     `json:"percent_visitors"`
	AvgCheck        float32 `json:"avg_check"`
	PrevAvgCheck    float32 `json:"prev_avg_check"`
	PercentAvgCheck int     `json:"percent_avg_check"`
}

type DaysOfTheWeek struct {
	Monday    Statistics `json:"monday"`
	Tuesday   Statistics `json:"tuesday"`
	Wednesday Statistics `json:"wednesday"`
	Thursday  Statistics `json:"thursday"`
	Friday    Statistics `json:"friday"`
	Saturday  Statistics `json:"saturday"`
	Sunday    Statistics `json:"sunday"`
}

type StatByHour struct {
	Sum      float32   `json:"sum"`
	Revenue  float32   `json:"revenue"`
	Profit   float32   `json:"profit"`
	Cost     float32   `json:"cost"`
	Checks   int       `json:"checks"`
	Visitors int       `json:"visitors"`
	AvgCheck float32   `json:"avg_check"`
	Hour     time.Time `json:"hour"`
}

type ABC struct {
	ID             int     `json:"id"`
	ItemID         int     `json:"item_id"`
	ItemName       string  `json:"item_name"`
	ItemType       string  `json:"item_type"`
	ItemCategory   int     `json:"item_category"`
	Sales          int     `json:"sales"`
	SalesPercent   float32 `json:"sales_percent"`
	SalesLetter    string  `json:"sales_letter"`
	Revenue        float32 `json:"revenue"`
	RevenuePercent float32 `json:"revenue_percent"`
	RevenueLetter  string  `json:"revenue_letter"`
	Profit         float32 `json:"profit"`
	ProfitPercent  float32 `json:"profit_percent"`
	ProfitLetter   string  `json:"profit_letter"`
	WorkerID       int     `json:"worker_id"`
	ShopID         int     `json:"shop_id"`
}
