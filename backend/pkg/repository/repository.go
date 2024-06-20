package repository

import (
	"database/sql"
	"time"
	"zebra/model"

	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Tovar interface {
	AddTovar(tovar *model.Tovar) error
	AddTovars(tovars []*model.Tovar) error
	GetAllTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error)
	GetTovar(id int, filter *model.Filter) (*model.TovarOutput, error)
	UpdateTovar(tovar *model.ReqTovar) error
	DeleteTovar(id int) error
	AddCategoryTovar(category *model.CategoryTovar) error
	GetAllCategoryTovar(filter *model.Filter) ([]*model.CategoryTovar, int64, error)
	GetCategoryTovar(id int) (*model.CategoryTovar, error)
	UpdateCategoryTovar(category *model.CategoryTovar) error
	DeleteCategoryTovar(id int) error
	AddTechCart(techCart *model.TechCart) error
	AddTechCarts(techCarts []*model.TechCart) error
	GetAllTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error)
	GetTechCart(id int, filter *model.Filter) (*model.TechCartOutput, error)
	UpdateTechCart(techCart *model.ReqTechCart, role string) error
	DeleteTechCart(id int) error
	GetTovarWithParams(sortParam, sklad, search string, category int) ([]*model.TovarOutput, error)
	GetTechCartWithParams(sortParam, sklad string, category int) ([]*model.TechCartOutput, error)
	GetTechCartNabor(id int) ([]*model.NaborOutput, error)
	GetEverything() ([]*model.ItemOutput, error)
	GetPureTovarByShopID(id, shopID int) (*model.Tovar, error)
	GetTovarsTovarIDByTovarsID(id int) (int, error)
	GetTechCartsTechCartIDByTechCartsID(id int) (int, error)
	GetDeletedTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error)
	GetDeletedTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error)
	RecreateTovar(tovar *model.Tovar) error
	RecreateTechCart(techCart *model.TechCart) error
	GetToAddTovar(filter *model.Filter) ([]*model.TovarMaster, int64, error)
	GetToAddTechCart(filter *model.Filter) ([]*model.TechCartMaster, int64, error)
	GetIdsOfShopsWhereTheTovarAlreadyExist(tovarID int) ([]int, error)
	GetIdsOfShopsWhereTheTechCartAlreadyExist(techCartID int) ([]int, error)
}

type Ingredient interface {
	AddIngredient(ingredient *model.Ingredient) error
	AddIngredients(ingredients []*model.Ingredient) error
	GetAllIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error)
	GetIngredient(id int, filter *model.Filter) (*model.IngredientOutput, error)
	UpdateIngredient(ingredient *model.ReqIngredient, shops []int) error
	DeleteIngredient(id int) error
	AddCategoryIngredient(category *model.CategoryIngredient) error
	GetAllCategoryIngredient(filter *model.Filter) ([]*model.CategoryIngredient, int64, error)
	GetCategoryIngredient(id int) (*model.CategoryIngredient, error)
	UpdateCategoryIngredient(category *model.CategoryIngredient) error
	DeleteCategoryIngredient(id int) error
	AddNabor(nabor *model.Nabor) (*model.Nabor, error)
	GetAllNabor(filter *model.Filter) ([]*model.NaborOutput, int64, error)
	GetNabor(id int) (*model.NaborOutput, error)
	UpdateNabor(nabor *model.Nabor) error
	DeleteNabor(id int) error
	GetTechCartByIngredientID(id int) ([]*model.TechCart, error)
	GetIngredientsForNabor(id int) ([]*model.Ingredient, error)
	GetPureIngredientByShopID(id, shopID int) (*model.Ingredient, error)
	GetIngredientsByTechCart(techCartID int) ([]*model.Ingredient, error)
	GetIngredientsIngredientIDByIngredientsID(id int) (int, error)
	GetDeletedIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error)
	RecreateIngredient(ingredient *model.Ingredient) error
	GetToAddIngredient(filter *model.Filter) ([]*model.IngredientMaster, int64, error)
	GetIdsOfShopsWhereTheIngredientAlreadyExist(ingredientID int) ([]int, error)
	GetIdsOfShopsWhereTheNaborAlreadyExist(naborID int) ([]int, error)
}

type Sklad interface {
	GetShopBySkladID(id int) (*model.Shop, error)
	CheckInventItems(items []*model.InventarizationItem, shopID int) error
	RecalculateNetCost() error
	RecalculateInventarizations() error
	AddSklad(sklad *model.Sklad) error
	GetAllSklad(filter *model.Filter) ([]*model.Sklad, int64, error)
	GetSklad(id int) (*model.Sklad, error)
	UpdateSklad(sklad *model.Sklad) error
	DeleteSklad(id int) error
	Ostatki(filter *model.Filter) ([]*model.Item, int64, error)
	AddToSklad(postavka *model.Postavka, shopID int, id int) (*model.Postavka, int, error)
	GetItems(filter *model.Filter) ([]*model.ItemOutput, error)
	GetAllPostavka(filter *model.Filter) ([]*model.PostavkaOutput, int64, error)
	GetSumOfPostavkaForPeriod(filter *model.Filter) (float32, error)
	GetPostavka(id int) (*model.PostavkaOutput, error)
	UpdatePostavka(postavka *model.Postavka) error
	DeletePostavka(id int) error
	RemoveFromSklad(spisanie *model.RemoveFromSklad) error
	RequestToRemove(request *model.RemoveFromSklad) error
	ConfirmToRemove(id int) error
	RejectToRemove(id int) error
	UpdateSpisanie(spisanie *model.RemoveFromSklad) error
	DeleteSpisanie(id int) error
	GetRemoved(filter *model.Filter) ([]*model.RemoveFromSkladResponse, int64, error)
	GetRemovedByID(id int) (*model.RemoveFromSkladResponse, error)
	GetTransactionFromPostavkaID(id int) (*model.Transaction, error)
	AddTransfer(transfer *model.Transfer) (*model.Transfer, error)
	GetAllTransfer(filter *model.Filter) ([]*model.TransferOutput, int64, error)
	GetTransfer(id int) (*model.TransferOutput, error)
	UpdateTransfer(transfer *model.Transfer) error
	DeleteTransfer(id int) error
	GetToCreateInventratization(inventarization *model.Inventarization) (*model.Inventarization, error)
	GetAllInventarization(filter *model.Filter) ([]*model.InventarizationResponse, int64, error)
	GetInventarization(id int) (*model.InventarizationResponse, error)
	UpdateInventarizationParams(inventarization *model.Inventarization) (*model.Inventarization, error)
	DeleteInventarization(id int) error
	DeleteInventarizationItem(id int) error
	UpdateInventarization(inventarization *model.Inventarization) (*model.Inventarization, error)
	UpdateInventarizationV2(inventarization *model.Inventarization) (*model.Inventarization, error)
	//GetShopBySkladID(id int) (*model.Shop, error)
	GetSkladByShopID(shopID int) (*model.Sklad, error)
	GetInventarizationDetailsIncome(id int) ([]*model.InventarizationDetailsIncome, error)
	GetInventarizationDetailsExpence(id int) ([]*model.InventarizationDetailsExpence, error)
	GetInventarizationDetailsSpisanie(id int) ([]*model.InventarizationDetailsSpisanie, error)
	DailyStatistic(shopID int) error
	GetTrafficReport(filter *model.Filter) ([]*model.TrafficReport, int64, error)
	RecalculateInventarization(invItems []*model.InventarizationItem) error
	RecalculateTrafficReport(*model.AsyncJob) error
	UpdateTrafficReportJob(*model.AsyncJob) error
	GetItemsForRecalculateTrafficReport() ([]*model.AsyncJob, error)
	GetInventarizationByID(id int) (*model.Inventarization, error)
	AddTrafficReportJob([]*model.AsyncJob) error
	ConcurrentRecalculationForDailyStatistics(items []*model.AsyncJob)
	CheckUnique(itemID, skladID int, itemType string, groupID int) error
	CreateInventarizationGroup(groupToAdd *model.InventarizationGroup) error
	GetAllInventarizationGroup(filter *model.Filter) ([]*model.InventarizationGroupResponse, error)
	GetInventarizationGroup(filter *model.Filter, id int) (*model.InventarizationGroupResponse, error)
	UpdateInventarizationGroup(group *model.InventarizationGroup) error
	DeleteInventarizationGroup(id int) error
	GetPureInventarizationGroup(id int) (*model.InventarizationGroup, error)
	GetPostavkaByTransferID(id int) (*model.Postavka, error)
	GetSpisanieByTransferID(id int) (*model.RemoveFromSklad, error)
	RecalculateDailyStatisticByDate(date time.Time, skladID int) error
	GetAllSklads() ([]*model.Sklad, error)
}

type Master interface {
	AddTovarMaster(tovar *model.TovarMaster) (*model.TovarMaster, error)
	AddIngredientMaster(ingredient *model.IngredientMaster) (*model.IngredientMaster, error)
	NormaliseTovars() error
	NormaliseIngredients() error
	NormaliseTechCarts() error
	NormaliseNabors() error
	AddTechCartMaster(techCart *model.TechCartMaster) (*model.TechCartMaster, error)
	AddNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	GetAllTovarMaster(filter *model.Filter) ([]*model.TovarMasterResponse, int64, error)
	GetTovarMaster(id int) (*model.TovarMasterResponse, error)
	GetAllIngredientMaster(filter *model.Filter) ([]*model.IngredientMasterResponse, int64, error)
	GetAllTechCartMaster(filter *model.Filter) ([]*model.TechCartMasterResponse, int64, error)
	GetIngredientMaster(id int) (*model.IngredientMasterResponse, error)
	GetTechCartMaster(id int) (*model.TechCartMasterResponse, error)
	GetAllNaborMaster(filter *model.Filter) ([]*model.NaborMasterOutput, int64, error)
	GetNaborMaster(id int) (*model.NaborMasterOutput, error)
	UpdateTovarMaster(tovar *model.TovarMaster) error
	UpdateIngredientMaster(ingredient *model.IngredientMaster) error
	UpdateTechCartMaster(techCart *model.TechCartMaster) error
	DeleteTovarMaster(id int) error
	DeleteIngredientMaster(id int) error
	DeleteTechCartMaster(id int) error
	DeleteNaborMaster(id int) error
	CreateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	UpdateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	UpdateTechCartsMaster(techCarts *model.TechCartMaster) error
	GetAllTovarMasterIds() ([]int, error)
	GetAllTechCartsMasterIds() ([]int, error)
	ConfirmTovarMaster(id int) (*model.TovarMaster, error)
	RejectTovarMaster(id int) error
	ConfirmTechCartMaster(id int) (*model.TechCartMaster, error)
	RejectTechCartMaster(id int) error
	ConfirmIngredientMaster(id int) (*model.IngredientMaster, error)
	RejectIngredientMaster(id int) error
	ConfirmNaborMaster(id int) (*model.NaborMaster, error)
	RejectNaborMaster(id int) error
}

type Mobile interface {
	CheckForRegister(email string) (bool, error)
	SendCode(email string) (string, error)
	Registrate(code string, req *model.ReqRegistrate) error
	GetClientByDeviceID(deviceID string) (*model.Client, error)
	UpdateCode(email, code string) error
	GetAllMobileUsers(filter *model.Filter) ([]*model.MobileUser, int64, error)
	GetMobileUser(id string, shopIDs []int) (*model.MobileUser, error)
	GetAllFeedbacks(filter *model.Filter) ([]*model.MobileUserFeedbackResponse, int64, error)
	GetFeedback(id int) (*model.MobileUserFeedbackResponse, error)
}

type Dealer interface {
	AddDealer(dealer *model.Dealer) error
	GetAllDealer(filter *model.Filter) ([]*model.Dealer, int64, error)
	GetDealerByID(id int) (*model.Dealer, error)
	UpdateDealer(dealer *model.Dealer) error
	DeleteDealer(id int) error
	GetAnyDealer() (*model.Dealer, error)
}

type Terminal interface {
	GetAllProducts(filter *model.Filter) ([]*model.Product, error)
}

type Shop interface {
	CreateShop(shop *model.Shop, products *model.ProductsShop) (*model.Shop, error)
	GetAllShop(filter *model.Filter) ([]*model.Shop, int64, error)
	GetShop(id int) (*model.Shop, error)
	UpdateShop(shop *model.Shop) error
	DeleteShop(id int) error
	GetShopBySchetID(id int) (*model.Shop, error)
	GetShopByCashSchetID(id int) (*model.Shop, error)
	GetAllShopsWithouParam() ([]*model.Shop, error)
	GetAllPureShops() ([]*model.Shop, error)
	GetRevenueByShopID(id int) (float32, error)
}

type Finance interface {
	AddSchet(schet *model.ReqSchet) (*model.Schet, error)
	GetAllSchet(filter *model.Filter) ([]*model.Schet, int64, error)
	GetSchetByID(id int) (*model.Schet, error)
	UpdateSchet(schet *model.Schet) error
	DeleteSchet(id int) error
}

type Check interface {
	UpdateToSend(check *model.SendToTis) error
	DeleteToSend(id int) error
	GetUnsendCheck() ([]*model.ErrorCheck, error)
	SaveError(fullError string, request string) error
	AddCheck(Check *model.Check) (*model.CheckResponse, error)
	UpdateCheck(check *model.Check) (*model.CheckResponse, error)
	CloseCheck(check *model.Check) (*model.Check, error)
	GetAllCheck(filter *model.Filter) ([]*model.Check, int64, error)
	GetAllWorkerCheck(filter *model.Filter) ([]*model.CheckResponse, int64, error)
	GetCheckByID(id int) (*model.CheckResponse, error)
	GetAllCheckView(page int) ([]*model.CheckView, int64, error)
	GetTisCheck(id int) (*model.ReqTisResponse, error)
	RemoveFromSklad(check *model.Check) (*model.Check, error)
	AddTag(tag *model.Tag) error
	GetAllTag(shopID int) ([]*model.Tag, error)
	GetTag(id int) (*model.Tag, error)
	UpdateTag(tag *model.Tag) error
	DeleteTag(id int) error
	GetModificatorsCost([]*model.Modificator) (float32, error)
	CalculateCheck(check *model.ReqCheck) (*model.Check, error)
	GetCheckByIdempotency(idempotency string) (*model.Check, error)
	SaveCheck(check *model.TisResponse) error
	DeleteCheck(id int) error
	GetTisToken(shopID int) (string, error)
	AddCheckToSend(check *model.SendToTis) error
	GetUnsendTisCheck() ([]*model.SendToTis, error)
	DeactivateCheck(id int) ([]*model.InventarizationItem, error)
	UpdateCheckLink(id int, link string) error
	IdempotencyCheck(keys *model.IdempotencyCheckArray, shopID int) (string, error)
	ConstructResponse(r *model.Check) *model.CheckResponse
	SaveFailedCheck(checks []*model.FailedCheck) error
	GetStoliki(shopID int) ([]*model.Stolik, error)
	GetFilledStoliki(shopID int) ([]*model.Stolik, error)
}
type External interface {
	SaveCheck(check *model.TisResponse) error
}

type User interface {
	AddUser(user *model.User) (int, error)
	GetUser(id int) (*model.User, error)
	GetCurrentOrders(id int) ([]*model.Check, error)
	AddFeedback(feedback *model.Feedback) error
	CheckCode(code string) (*model.UserQR, error)
	SetCode(id int, code string) (*model.UserQR, error)
	GetUserCode(id int) (*model.UserQR, error)
	GetUserByCode(code string) (*model.User, error)
	CleanExpiredCode() error
	GetWorkerByUsername(username string) (*model.Worker, error)
	CreateWorker(worker *model.Worker) (*model.Worker, error)
	UpdateWorker(worker *model.Worker) error
}

type Worker interface {
	GetAllWorkers(filter *model.Filter) ([]*model.Worker, int64, error)
	GetWorker(id int) (*model.Worker, error)
	UpdateWorker(worker *model.Worker) (*model.Worker, error)
	DeleteWorker(id int) error
}

type Statistics interface {
	GetWorkersStat(filter *model.Filter) ([]*model.WorkerStat, int64, error)
	TodayStatistics(time time.Time, filter *model.Filter) (*model.TodayStatistics, error)
	EveryDayStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error)
	EveryWeekStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error)
	EveryMonthStatistics(time time.Time, filter *model.Filter) ([]*model.Statistics, error)
	Payments(filter *model.Filter) ([]*model.Payment, int64, error)
	DaysOfTheWeek(filter *model.Filter) (*model.DaysOfTheWeek, error)
	StatByHour(filter *model.Filter) ([]*model.StatByHour, error)
	ABC(filter *model.Filter) ([]*model.ABC, error)
	TopSales(filter *model.Filter) ([]*model.ItemOutput, error)
}

type Transaction interface {
	OpenShift(id int, shift *model.Shift) (*model.Shift, error)
	CloseShift(id int, shift *model.Shift) error
	UpdateShift(shift *model.Shift) error
	CreateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	GetLastShift(shopID int) (*model.Shift, error)
	CountOpenCheck(id int, time time.Time) (int, error)
	GetAllShifts(filter *model.Filter) ([]*model.Shift, int64, error)
	GetShiftByTime(timeStamp time.Time, shopID int) (*model.Shift, error)
	GetShiftByID(id int) (*model.ShiftTransaction, error)
	GetShiftByID2(id int) (*model.Shift, error)
	GetShiftByShopId(id int, workerID int) (*model.CurrentShift, error)
	CreatePostavkaTransaction(transactionPostavka *model.TransactionPostavka) error
	GetAllTransaction(filter *model.Filter) ([]*model.TransactionResponse, int64, error)
	GetTransactionByID(id int) (*model.TransactionResponse, error)
	GetTransactionsByID(id int) ([]*model.Transaction, error)
	UpdateTransaction(transaction *model.Transaction) error
	DeleteTransaction(id int) error
	RecalculateShift(id int) error
	CheckForBlockedShop(shiftID int) (bool, error)
}

type Repository struct {
	Terminal
	Tovar
	Ingredient
	Sklad
	Dealer
	Finance
	Check
	External
	User
	Worker
	Statistics
	Transaction
	Shop
	Master
	Mobile
}

func NewRepository(db *sql.DB, gormDB *gorm.DB) *Repository {
	return &Repository{
		Terminal:    NewTerminalDB(db, gormDB),
		Tovar:       NewTovarDB(db, gormDB),
		Ingredient:  NewIngredientDB(db, gormDB),
		Sklad:       NewSkladDB(db, gormDB),
		Dealer:      NewDealerDB(db, gormDB),
		Finance:     NewFinanceDB(db, gormDB),
		Check:       NewCheckDB(db, gormDB),
		External:    NewExternalDB(db, gormDB),
		User:        NewUserDB(db, gormDB),
		Worker:      NewWorkerDB(db, gormDB),
		Statistics:  NewStatisticsDB(db, gormDB),
		Transaction: NewTransactionDB(db, gormDB),
		Shop:        NewShopDB(db, gormDB),
		Master:      NewMasterDB(db, gormDB),
		Mobile:      NewMobileDB(db, gormDB),
	}
}
