package model

type Product struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	CategoryName string         `json:"categoryName"`
	Category     int            `json:"category"`
	Type         string         `json:"type"`
	Image        string         `json:"image"`
	Price        float32        `json:"price" gorm:"column:price;not null"`
	Discount     bool           `json:"discount" gorm:"default:false"`
	Nabor        []*NaborOutput `json:"nabor"`
}

type Terminal struct {
	Categories   []*CategoryProduct `json:"categories"`
	MainDisplay  []*Product         `json:"mainDisplay"`
	CurrentShift *CurrentShift      `json:"current_shift"`
}

type CategoryProduct struct {
	ID       int        `json:"id"`
	Category string     `json:"category"`
	Image    string     `json:"image"`
	Products []*Product `json:"products"`
}
