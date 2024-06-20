package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

type Sklad struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Address         string             `json:"address"`
	Sum             float32            `json:"sum"`
	ShopID          int                `json:"shop_id"`
	Deleted         bool               `json:"deleted" gorm:"default:false"`
	Postavka        []*Postavka        `json:"postavka" gorm:"ForeignKey:SkladID"`
	RemoveFromSklad []*RemoveFromSklad `json:"remove_from_sklad" gorm:"ForeignKey:SkladID"`
}

func (p *Sklad) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if p.Name == "" || p.Address == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type SkladTovar struct {
	SkladID  int     `json:"sklad_id"`
	TovarID  int     `json:"tovar_id"`
	Quantity float32 `json:"quantity"`
	Cost     float32 `json:"cost"`
}

type SkladIngredient struct {
	SkladID      int     `json:"sklad_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float32 `json:"quantity"`
	Cost         float32 `json:"cost"`
}

type Postavka struct {
	ID                  int                 `json:"id"`
	DealerID            int                 `json:"dealer_id"`
	SkladID             int                 `json:"sklad_id"`
	SchetID             int                 `json:"schet_id"`
	Time                time.Time           `json:"time"`
	Sum                 float32             `json:"sum"`
	Risky               bool                `json:"risky"`
	Deleted             bool                `json:"deleted" gorm:"default:false"`
	Type                string              `json:"transfer"`
	TransferID          int                 `json:"transfer_id"`
	Items               []ItemPostavka      `json:"items" gorm:"ForeignKey:PostavkaID"`
	TransactionPostavka TransactionPostavka `json:"transaction_postavka" gorm:"ForeignKey:PostavkaID"`
	Comment             string              `json:"comment"`
}

type PostavkaAll struct {
	ID            int              `json:"postavka_id"`
	PostavkaItems PostavkaItemsArr `json:"postavka_items"`
}

type PostavkaItems struct {
	PostavkaInfo []*PostavkaInfo `json:"postavka_items" gorm:"ForeignKey:PostavkaID"`
}

type PostavkaItemsArr []*PostavkaItems

func (p *PostavkaItemsArr) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var items []*PostavkaInfo
	err := json.Unmarshal(value.([]byte), &items)
	if err != nil {
		return err
	}
	var postavka2Items []*PostavkaItems
	for _, item := range items {
		postavka2Items = append(postavka2Items, &PostavkaItems{PostavkaInfo: []*PostavkaInfo{item}})
	}
	*p = postavka2Items
	return nil
}

type PostavkaInfo struct {
	Sum         float32   `json:"sum"`
	Dealer      string    `json:"dealer"`
	Sklad       string    `json:"sklad"`
	Schet       string    `json:"schet"`
	Time        time.Time `json:"time"`
	ItemID      int       `json:"item_id"`
	PostavkaID  int       `json:"postavka_id"`
	Type        string    `json:"type"`
	Quantity    float32   `json:"quantity"`
	Cost        float32   `json:"cost"`
	Name        string    `json:"name"`
	Measurement string    `json:"measurement"`
	Category    string    `json:"category"`
	Risky       bool      `json:"risky"`
	Comment     string    `json:"comment"`
	Deleted     bool      `json:"deleted"`
}

type GlobalPostavka struct {
	Sum      float32           `json:"sum"`
	Postavka []*PostavkaOutput `json:"postavka"`
}

type PostavkaOutput struct {
	ID       int                   `json:"id"`
	Dealer   string                `json:"dealer"`
	DealerID int                   `json:"dealer_id"`
	Sklad    string                `json:"sklad"`
	SkladID  int                   `json:"sklad_id"`
	Schet    string                `json:"schet"`
	SchetID  int                   `json:"schet_id"`
	Time     time.Time             `json:"time"`
	Sum      float32               `json:"sum"`
	Risky    bool                  `json:"risky"`
	Comment  string                `json:"comment"`
	Deleted  bool                  `json:"deleted"`
	Items    []*ItemPostavkaOutput `json:"items" gorm:"ForeignKey:PostavkaID"`
}

type ItemPostavkaOutput struct {
	ID          int     `json:"id"`
	ItemID      int     `json:"item_id"`
	PostavkaID  int     `json:"postavka_id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Measurement string  `json:"measurement"`
	Quantity    float32 `json:"quantity"`
	Cost        float32 `json:"cost"`
	Risky       bool    `json:"risky"`
	Details     string  `json:"details"`
	Deleted     bool    `json:"deleted"`
}

type NetCost struct {
	Cost     float32 `json:"cost"`
	Quantity float32 `json:"quantity"`
}

type ChangeNetCost struct {
	ItemID  int    `json:"item_id"`
	Type    string `json:"type"`
	SkladID int    `json:"sklad_id"`
}

type ItemPostavka struct {
	ID         int     `json:"id" gorm:"primary_key"`
	ItemID     int     `json:"item_id"`
	PostavkaID int     `json:"postavka_id"`
	Type       string  `json:"type"`
	Quantity   float32 `json:"quantity"`
	Cost       float32 `json:"cost"`
	Risky      bool    `json:"risky"`
	Details    string  `json:"details"`
	Deleted    bool    `json:"deleted" gorm:"default:false"`
}

func (p *ItemPostavka) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.ItemID <= 0 || p.Quantity <= 0 || p.Cost <= 0 {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

func (p *Postavka) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if len(p.Items) <= 0 {
		return errors.New("bad request | fill fields properly")
	}
	p.DealerID = 1

	if p.Time.IsZero() {
		p.Time = time.Now()
	}

	for _, item := range p.Items {
		if item.ItemID <= 0 || item.Quantity <= 0 || item.Cost <= 0 {
			return errors.New("bad request | fill fields properly")
		}
	}

	return nil
}

type Item struct {
	ItemID      int     `json:"id"`
	SkladID     int     `json:"skladID"`
	Name        string  `json:"name"`
	SkladName   string  `json:"skladName"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Quantity    float32 `json:"quantity"`
	Cost        float32 `json:"cost"`
	Measurement string  `json:"measurement"`
	Sum         float32 `json:"sum"`
}

type ItemOutput struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Measure          string    `json:"measure"`
	Cost             float32   `json:"cost"`
	SkladID          int       `json:"sklad_id"`
	LastPostavkaTime time.Time `json:"last_postavka_time"`
	LastPostavkaCost float32   `json:"last_postavka_cost"`
	Type             string    `json:"type"`
	Sales            float32   `json:"sales"`
}

type GlobalSpisanie struct {
	Sum             float32                    `json:"sum"`
	RemoveFromSklad []*RemoveFromSkladResponse `json:"remove_from_sklad"`
}

type RemoveFromSkladResponse struct {
	ID       int                            `json:"id" gorm:"primary_key"`
	Sklad    string                         `json:"sklad"`
	SkladID  int                            `json:"sklad_id"`
	WorkerID int                            `json:"worker_id"`
	Reason   string                         `json:"reason"`
	Comment  string                         `json:"comment"`
	Cost     float32                        `json:"cost"`
	Time     time.Time                      `json:"time"`
	Status   string                         `json:"status"`
	Items    []*RemoveFromSkladItemResponse `json:"items" gorm:"ForeignKey:ItemID"`
}

type RemoveFromSkladOutput struct {
	ID   int                      `json:"id"`
	Info ItemsRemovedFromSkladArr `json:"remove_from_sklads_info"`
}

type ItemsRemovedFromSklad struct {
	Items []*RemoveFromSkladInfo `json:"remove_from_sklads_info" gorm:"ForeignKey:Remove_from_sklads_infoID"`
}

func (r *ItemsRemovedFromSkladArr) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var items []*RemoveFromSkladInfo
	err := json.Unmarshal(value.([]byte), &items)
	if err != nil {
		return err
	}
	var postavka2Items []*ItemsRemovedFromSklad
	for _, item := range items {
		postavka2Items = append(postavka2Items, &ItemsRemovedFromSklad{Items: []*RemoveFromSkladInfo{item}})
	}
	*r = postavka2Items
	return nil
}

type RemoveFromSkladInfo struct {
	Sklad    string    `json:"sklad"`
	WorkerID int       `json:"worker_id"`
	Reason   string    `json:"reason"`
	Comment  string    `json:"comment"`
	Cost     float32   `json:"cost"`
	Time     time.Time `json:"time"`
	Status   string    `json:"status"`
	ItemID   int       `json:"item_id"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Measure  string    `json:"measure"`
	Quantity float32   `json:"quantity"`
	ItemCost float32   `json:"item_cost"`
	Details  string    `json:"details"`
}

type ItemsRemovedFromSkladArr []*ItemsRemovedFromSklad

type RemoveFromSkladItemResponse struct {
	ID       int     `json:"id" gorm:"primary_key"`
	ItemID   int     `json:"item_id"`
	Name     string  `json:"name"`
	Measure  string  `json:"measure"`
	Type     string  `json:"type"`
	Quantity float32 `json:"quantity"`
	Cost     float32 `json:"cost"`
	Details  string  `json:"details"`
}

type RemoveFromSklad struct {
	ID         int                    `json:"id" gorm:"primary_key"`
	SkladID    int                    `json:"sklad_id"`
	WorkerID   int                    `json:"worker_id"`
	Reason     string                 `json:"reason"`
	Comment    string                 `json:"comment"`
	Cost       float32                `json:"cost"`
	Time       time.Time              `json:"time"`
	Status     string                 `json:"status"`
	Type       string                 `json:"transfer" gorm:"default:false"`
	TransferID int                    `json:"transfer_id"`
	Deleted    bool                   `json:"deleted" gorm:"default:false"`
	Items      []*RemoveFromSkladItem `json:"items" gorm:"ForeignKey:RemoveID"`
}

func (p *RemoveFromSklad) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Time.IsZero() {
		p.Time = time.Now()
	}

	if len(p.Items) <= 0 {
		return errors.New("bad request | fill fields properly")
	}
	for _, item := range p.Items {
		if item.ItemID <= 0 || item.Quantity <= 0 {
			return errors.New("bad request | fill fields properly")
		}
	}

	return nil
}

type RemoveFromSkladItem struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RemoveID       int     `json:"remove_id"`
	ItemID         int     `json:"item_id"`
	SkladID        int     `json:"sklad_id"`
	PartOfTechCart bool    `json:"part_of_tech_cart" gorm:"default:false"`
	Type           string  `json:"type"`
	Quantity       float32 `json:"quantity"`
	Cost           float32 `json:"cost"`
	Details        string  `json:"details"`
}

type Inventarization struct {
	ID                   int                    `json:"id" gorm:"primary_key"`
	SkladID              int                    `json:"sklad_id"`
	Time                 time.Time              `json:"time"`
	Type                 string                 `json:"type"`
	Result               float32                `json:"result"`
	Status               string                 `json:"status"`
	LoadingStatus        string                 `json:"loading_status"`
	Deleted              bool                   `json:"deleted"`
	InventarizationItems []*InventarizationItem `json:"items" gorm:"ForeignKey:InventarizationID"`
}

type InventarizationForNetCost struct {
	Time         time.Time `json:"time"`
	FactQuantity float32   `json:"fact_quantity"`
	Cost         float32   `json:"cost"`
}

func (p *Inventarization) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if p.Time.IsZero() {
		p.Time = time.Now()
	}

	if p.Type != utils.TypePartial && p.Type != utils.TypeFull && p.Type != utils.TypeFullPartial {
		return errors.New("bad request | fill fields properly")
	}
	if p.Status != utils.StatusClosed && p.Status != utils.StatusOpened {
		return errors.New("bad request | fill fields properly")
	}
	for _, item := range p.InventarizationItems {
		if item.ItemID <= 0 || item.Type == "" {
			return errors.New("bad request | fill fields properly")
		}
	}
	return nil
}

type InventarizationItem struct {
	ID                int       `json:"id" gorm:"primary_key"`
	InventarizationID int       `json:"inventarization_id"`
	ItemID            int       `json:"item_id"`
	SkladID           int       `json:"sklad_id"`
	Status            string    `json:"status"`
	Time              time.Time `json:"time"`
	Type              string    `json:"type"`
	StartQuantity     float32   `json:"start_quantity"`
	Expenses          float32   `json:"expenses"`
	Income            float32   `json:"income"`
	Removed           float32   `json:"removed"`
	RemovedSum        float32   `json:"removed_sum"`
	PlanQuantity      float32   `json:"plan_quantity"`
	FactQuantity      float32   `json:"fact_quantity"`
	Difference        float32   `json:"difference"`
	DifferenceSum     float32   `json:"difference_sum"`
	Cost              float32   `json:"cost"`
	LoadingStatus     string    `json:"loading_status"`
	BeforeTime        time.Time `json:"before_time"`
	NeedToRecalculate bool      `json:"need_to_recalculate"`
	GroupID           int       `json:"group_id"`
	IsVisible         bool      `json:"is_visible"`
}

func (invItem *InventarizationItem) HasChanged(oldItem *InventarizationItem) bool {
	return !oldItem.Time.Equal(invItem.Time) || oldItem.FactQuantity != invItem.FactQuantity || oldItem.Status != invItem.Status || oldItem.NeedToRecalculate || oldItem.GroupID != invItem.GroupID
}

type InventarizationResponse struct {
	ID                   int                            `json:"id" gorm:"primary_key"`
	Sklad                string                         `json:"sklad"`
	SkladID              int                            `json:"sklad_id"`
	Time                 time.Time                      `json:"time"`
	Type                 string                         `json:"type"`
	Result               float32                        `json:"result"`
	LoadingStatus        string                         `json:"loading_status"`
	Status               string                         `json:"status"`
	InventarizationItems []*InventarizationItemResponse `json:"items" gorm:"ForeignKey:InventarizationID"`
}

type InventarizationItemResponse struct {
	ID                int       `json:"id" gorm:"primary_key"`
	InventarizationID int       `json:"inventarization_id"`
	ItemID            int       `json:"item_id"`
	SkladID           int       `json:"sklad_id"`
	ItemName          string    `json:"item_name"`
	SkladName         string    `json:"sklad_name"`
	Status            string    `json:"status"`
	Time              time.Time `json:"time"`
	BeforeTime        time.Time `json:"before_time"`
	Type              string    `json:"type"`
	StartQuantity     float32   `json:"start_quantity"`
	Expenses          float32   `json:"expenses"`
	Income            float32   `json:"income"`
	Removed           float32   `json:"removed"`
	RemovedSum        float32   `json:"removed_sum"`
	PlanQuantity      float32   `json:"plan_quantity"`
	FactQuantity      float32   `json:"fact_quantity"`
	Difference        float32   `json:"difference"`
	DifferenceSum     float32   `json:"difference_sum"`
	Measure           string    `json:"measure"`
	Cost              float32   `json:"cost"`
	LoadingStatus     string    `json:"loading_status"`
	IsVisible         bool      `json:"is_visible"`
	GroupID           int       `json:"group_id"`
}
type Transfer struct {
	ID            int             `json:"id"`
	Time          time.Time       `json:"time"`
	FromSklad     int             `json:"from_sklad"`
	ToSklad       int             `json:"to_sklad"`
	Worker        int             `json:"worker"`
	Comment       string          `json:"comment"`
	ItemTransfers []*ItemTransfer `json:"item_transfers" gorm:"ForeignKey:TransferID"`
	Sum           float32         `json:"sum"`
	Less          bool            `json:"less" gorm:"default:false"`
	Deleted       bool            `json:"deleted" gorm:"default:false"`
}

type ItemTransfer struct {
	ID         int     `json:"id"`
	ItemID     int     `json:"item_ID"`
	TransferID int     `json:"transfer_id"`
	Type       string  `json:"type"`
	Quantity   float32 `json:"quantity"`
	Sum        float32 `json:"sum"`
}

type TransferOutput struct {
	ID            int                   `json:"id"`
	Time          time.Time             `json:"time"`
	FromSklad     int                   `json:"from_sklad"`
	FromSkladName string                `json:"from_sklad_name"`
	ToSklad       int                   `json:"to_sklad"`
	ToSkladName   string                `json:"to_sklad_name"`
	Worker        int                   `json:"worker"`
	WorkerName    string                `json:"worker_name"`
	Comment       string                `json:"comment"`
	ItemTransfers []*ItemTransferOutput `json:"item_transfers" gorm:"ForeignKey:TransferID"`
	Sum           float32               `json:"sum"`
}

type ItemTransferOutput struct {
	ID          int     `json:"id"`
	ItemID      int     `json:"item_ID"`
	TransferID  int     `json:"transfer_id"`
	Name        string  `json:"item_name"`
	Type        string  `json:"type"`
	Quantity    float32 `json:"quantity"`
	Measurement string  `json:"measurement"`
	Category    string  `json:"category"`
	Cost        float32 `json:"cost"`
	Sum         float32 `json:"sum"`
}

func (t *Transfer) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&t); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if t.FromSklad <= 0 || t.ToSklad <= 0 || len(t.ItemTransfers) <= 0 {
		return errors.New("bad request | fill fields properly")
	}

	if t.Time.IsZero() {
		t.Time = time.Now()
	}

	for _, item := range t.ItemTransfers {
		if item.ItemID <= 0 || item.Quantity <= 0 {
			return errors.New("bad request | fill fields properly")
		}
	}

	return nil
}

type ExpenceTovar struct {
	ID           int       `json:"id" gorm:"primary_key"`
	SkladID      int       `json:"sklad_id"`
	TovarID      int       `json:"tovar_id"`
	Quantity     float32   `json:"quantity"`
	Cost         float32   `json:"cost"`
	Time         time.Time `json:"time"`
	CheckTovarID int       `json:"check_tovar_id"`
	Status       string    `json:"status"`
}
type ExpenceIngredient struct {
	ID              int       `json:"id" gorm:"primary_key"`
	SkladID         int       `json:"sklad_id"`
	IngredientID    int       `json:"ingerdient_id"`
	Quantity        float32   `json:"quantity"`
	Cost            float32   `json:"cost"`
	Time            time.Time `json:"time"`
	Type            string    `json:"type"`
	Price           float32   `json:"price"`
	CheckTechCartID int       `json:"check_tech_cart_id"`
	Status          string    `json:"status"`
}

type DailyStatistic struct {
	Date            time.Time `json:"date"`
	ID              int       `json:"id"`
	ItemID          int       `json:"item_id"`
	ItemName        string    `json:"item_name"`
	CheckCost       float32   `json:"check_cost"`
	CheckPrice      float32   `json:"check_price"`
	SkladID         int       `json:"sklad_id"`
	ShopID          int       `json:"shop_id"`
	Type            string    `json:"type"`
	Postavka        float32   `json:"postavka"`
	PostavkaCost    float32   `json:"postavka_cost"`
	Inventarization float32   `json:"inventarization"`
	Transfer        float32   `json:"transfer"`
	Sales           float32   `json:"sales"`
	RemoveFromSklad float32   `json:"remove_from_sklad"`
	Cost            float32   `json:"cost"`
	Quantity        float32   `json:"quantity"`
}

type GlobalTrafficReport struct {
	InitialSum     float32          `json:"initial_sum"`
	FinalSum       float32          `json:"final_sum"`
	TrafficReports []*TrafficReport `json:"traffic_reports"`
}

type TrafficReport struct {
	ItemID          int     `json:"item_id"`
	ItemName        string  `json:"item_name"`
	Measure         string  `json:"measure"`
	SkladID         int     `json:"sklad_id"`
	Type            string  `json:"type"`
	InitialOstatki  float32 `json:"initial_ostatki"`
	InitialNetCost  float32 `json:"initial_netCost"`
	InitialSum      float32 `json:"initial_sum"`
	Income          float32 `json:"income"`
	Consumption     float32 `json:"consumption"`
	FinalOstatki    float32 `json:"final_ostatki"`
	FinalNetCost    float32 `json:"final_netCost"`
	FinalSum        float32 `json:"final_sum"`
	Postavka        float32 `json:"postavka"`
	PostavkaCost    float32 `json:"postavka_cost"`
	Inventarization float32 `json:"inventarization"`
	Sales           float32 `json:"sales"`
	RemoveFromSklad float32 `json:"remove_from_sklad"`
	Transfer        float32 `json:"transfer"`
}

type InventarizationDetailsIncome struct {
	Time        time.Time `json:"time"`
	Quantity    float32   `json:"quantity"`
	Name        string    `json:"name"`
	Cost        float32   `json:"cost"`
	Sum         float32   `json:"sum"`
	Measurement string    `json:"measurement"`
}

type InventarizationDetailsExpence struct {
	Time        time.Time `json:"time"`
	Name        string    `json:"name"`
	Quantity    float32   `json:"quantity"`
	Measurement string    `json:"measurement"`
}

type InventarizationDetailsSpisanie struct {
	Time        time.Time `json:"time"`
	Quantity    float32   `json:"quantity"`
	Name        string    `json:"name"`
	Cost        float32   `json:"cost"`
	Sum         float32   `json:"sum"`
	Measurement string    `json:"measurement"`
}

type AsyncJob struct {
	ID         int       `json:"id"`
	ItemID     int       `json:"item_id"`
	ItemType   string    `json:"item_type"`
	SkladID    int       `json:"sklad_id"`
	ShopID     int       `json:"shop_id"`
	TimeStamp  time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Exception  string    `json:"exception"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
}

type InventarizationGroup struct {
	ID      int                         `json:"id" gorm:"primary_key"`
	Name    string                      `json:"name"`
	Measure string                      `json:"measure"`
	Type    string                      `json:"type"`
	SkladID int                         `json:"sklad_id"`
	Items   []*InventarizationGroupItem `json:"items"  gorm:"ForeignKey:GroupID"`
}

type InventarizationGroupItem struct {
	ID      int `json:"id"`
	ItemID  int `json:"item_id"`
	SkladID int `json:"sklad_id"`
	GroupID int `json:"group_id"`
}

type InventarizationGroupResponse struct {
	ID        int                                 `json:"id"`
	Name      string                              `json:"name"`
	Measure   string                              `json:"measure"`
	Type      string                              `json:"type"`
	SkladID   int                                 `json:"sklad_id"`
	SkladName string                              `json:"sklad_name"`
	Items     []*InventarizationGroupItemResponse `json:"items" gorm:"-"`
}

type InventarizationGroupItemResponse struct {
	ID      int    `json:"id"`
	ItemID  int    `json:"item_id"`
	GroupID int    `json:"group_id"`
	SkladID int    `json:"sklad_id"`
	Name    string `json:"name"`
}
