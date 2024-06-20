package service

import (
	"zebra/model"
	"zebra/pkg/repository"
)

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
	GetTovarWithParams(sort, sklad, search string, category int) ([]*model.TovarOutput, error)
	GetTechCartWithParams(sort, sklad string, category int) ([]*model.TechCartOutput, error)
	GetTechCartNabor(id int) ([]*model.NaborOutput, error)
	GetEverything() ([]*model.ItemOutput, error)
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

type Terminal interface {
	GetAllProducts(filter *model.Filter) (*model.Terminal, error)
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
	AddNabors(nabors []*model.Nabor) ([]*model.Nabor, error)
	GetTechCartByIngredientID(id int) ([]*model.TechCart, error)
	GetIngredientsIngredientIDByIngredientsID(id int) (int, error)
	GetDeletedIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error)
	RecreateIngredient(ingredient *model.Ingredient) error
	GetToAddIngredient(filter *model.Filter) ([]*model.IngredientMaster, int64, error)
	GetIdsOfShopsWhereTheIngredientAlreadyExist(ingredientID int) ([]int, error)
	GetIdsOfShopsWhereTheNaborAlreadyExist(naborID int) ([]int, error)
}

type Shop interface {
	CreateShop(shop *model.ReqShop, products *model.ProductsShop) (*model.Shop, error)
	GetAllShops(filter *model.Filter) ([]*model.Shop, int64, error)
	GetShopByID(id int) (*model.Shop, error)
	GetShopBySchetID(id int) (*model.Shop, error)
	GetShopByCashSchetID(id int) (*model.Shop, error)
	UpdateShop(shop *model.ReqShop) error
	DeleteShop(id int) error
}

type Sklad interface {
	RecalculateNetCost() error
	RecalculateInventarization() error
	AddSklad(sklad *model.Sklad) error
	GetAllSklad(filter *model.Filter) ([]*model.Sklad, int64, error)
	GetSklad(id int) (*model.Sklad, error)
	UpdateSklad(sklad *model.Sklad) error
	DeleteSklad(id int) error
	Ostatki(filter *model.Filter) ([]*model.Item, int64, error)
	AddToSklad(postavka *model.Postavka, shopID int, id int) (*model.Postavka, error)
	GetItems(filter *model.Filter) ([]*model.ItemOutput, error)
	GetAllPostavka(filter *model.Filter) (*model.GlobalPostavka, int64, error)
	GetPostavka(id int) (*model.PostavkaOutput, error)
	UpdatePostavka(postavka *model.Postavka) error
	DeletePostavka(id int) error
	RemoveFromSklad(spisanie *model.RemoveFromSklad) error
	RequestToRemove(request *model.RemoveFromSklad) error
	ConfirmToRemove(id int) error
	RejectToRemove(id int) error
	UpdateSpisanie(spisanie *model.RemoveFromSklad) error
	DeleteSpisanie(id int) error
	GetRemoved(filter *model.Filter) (*model.GlobalSpisanie, int64, error)
	GetRemovedByID(id int) (*model.RemoveFromSkladResponse, error)
	AddTransfer(transfer *model.Transfer) error
	GetAllTransfer(filter *model.Filter) ([]*model.TransferOutput, int64, error)
	GetTransfer(id int) (*model.TransferOutput, error)
	UpdateTransfer(transfer *model.Transfer) error
	DeleteTransfer(id int) error

	GetToCreateInventratization(inventarization *model.Inventarization) (*model.Inventarization, error)
	OpenInventarization(inventarization *model.Inventarization) (*model.Inventarization, error)
	UpdateInventarization(inventarization *model.Inventarization) (*model.Inventarization, error)
	UpdateInventarizationV2(inventarization *model.Inventarization) (*model.Inventarization, error)
	UpdateInventarizationParams(inventarization *model.Inventarization) (*model.Inventarization, error)
	GetAllInventarization(filter *model.Filter) ([]*model.InventarizationResponse, int64, error)
	GetInventarization(id int) (*model.InventarizationResponse, error)
	DeleteInventarization(id int) error
	DeleteInventarizationItem(id int) error

	//GetShopBySkladID(id int) (*model.Shop, error)
	GetSkladByShopID(shopID int) (*model.Sklad, error)
	GetInventarizationDetailsIncome(id int) ([]*model.InventarizationDetailsIncome, error)
	GetInventarizationDetailsExpence(id int) ([]*model.InventarizationDetailsExpence, error)
	GetInventarizationDetailsSpisanie(id int) ([]*model.InventarizationDetailsSpisanie, error)
	GetTrafficReport(filter *model.Filter) (*model.GlobalTrafficReport, int64, error)
	DailyStatistic() error
	RecalculateTrafficReport() error
	CreateInventarizationGroup(group *model.ReqInventarizationGroup) error
	GetAllInventarizationGroup(filter *model.Filter) ([]*model.InventarizationGroupResponse, error)
	GetInventarizationGroup(filter *model.Filter, id int) (*model.InventarizationGroupResponse, error)
	UpdateInventarizationGroup(group *model.ReqInventarizationGroup) error
	DeleteInventarizationGroup(id int) error
	RecalculateDailyStatistic() error
}

type Dealer interface {
	AddDealer(dealer *model.Dealer) error
	GetAllDealer(filter *model.Filter) ([]*model.Dealer, int64, error)
	GetDealer(id int) (*model.Dealer, error)
	UpdateDealer(dealer *model.Dealer) error
	DeleteDealer(id int) error
}

type Finance interface {
	AddSchet(schet *model.ReqSchet) (*model.Schet, error)
	GetAllSchet(filter *model.Filter) ([]*model.Schet, int64, error)
	GetSchet(id int) (*model.Schet, error)
	UpdateSchet(schet *model.Schet) error
	DeleteSchet(id int) error
}

type Check interface {
	SaveError(fullError string, request string)
	SendToTis() error
	AddCheck(Check *model.ReqCheck) (*model.CheckResponse, error)
	UpdateCheck(check *model.ReqCheck) (*model.CheckResponse, error)
	GetAllCheck(filter *model.Filter) (*model.GlobalCheck, int64, error)
	GetAllWorkerCheck(filter *model.Filter) ([]*model.CheckResponse, int64, error)
	GetCheckByID(id int) (*model.CheckResponse, error)
	GetAllCheckView(page int) ([]*model.CheckView, int64, error)
	GetCheckByIDForPrinter(id int) (*model.CheckPrint, error)
	AddTag(tag *model.Tag) error
	GetAllTag(shopID int) ([]*model.Tag, error)
	GetTag(id int) (*model.Tag, error)
	UpdateTag(tag *model.Tag) error
	DeleteTag(id int) error
	GetCheckByIdempotency(idempotency string) (*model.Check, error)
	CreateTisBody(check *model.CheckResponse, tisType int) (*model.TisBody, error)
	DeleteCheck(id int) error
	DeactivateCheck(id int) error
	IdempotencyCheck(keys *model.IdempotencyCheckArray, shopID int) error
	SaveFailedCheck(checks []*model.FailedCheck) error
	GetStoliki(shopID int) ([]*model.Stolik, error)
	GetFilledStoliki(shopID int) ([]*model.Stolik, error)
}

type External interface {
	SaveCheck(check *model.ReqTisResponse) error
}

type User interface {
	AddUser(user *model.UserRequest) (int, error)
	GetUser(id int) (*model.User, error)
	GetCurrentOrders(id int) ([]*model.Check, error)
	AddFeedback(feedback *model.FeedbackRequest) error
	GenCode(id int) (string, int64, error)
	GetUserByCode(code string) (*model.User, error)
	GetWorkerByUsername(username string) (*model.Worker, error)
	CreateWorker(worker *model.ReqWorkerRegistration) error
	ParseToken(accessToken string) (int, string, []int, int, error)
	GenerateToken(userID int, role string, shops []int, bindShop int) (string, error)
}

type Worker interface {
	GetAllWorkers(filter *model.Filter) ([]*model.Worker, int64, error)
	GetWorker(id int) (*model.Worker, error)
	UpdateWorker(worker *model.Worker) (*model.Worker, error)
	DeleteWorker(id int) error
}

type Statistics interface {
	GetWorkersStat(filter *model.Filter) ([]*model.WorkerStat, int64, error)
	TodayStatistics(filter *model.Filter) (*model.TodayStatistics, error)
	EveryDayStatistics(filter *model.Filter) (*model.TotalStatistics, error)
	EveryWeekStatistics(filter *model.Filter) (*model.TotalStatistics, error)
	EveryMonthStatistics(filter *model.Filter) (*model.TotalStatistics, error)
	Payments(filter *model.Filter) (*model.GlobalPayment, int64, error)
	DaysOfTheWeek(filter *model.Filter) (*model.DaysOfTheWeek, error)
	StatByHour(filter *model.Filter) ([]*model.StatByHour, error)
	ABC(filter *model.Filter) ([]*model.ABC, error)
	TopSales(filter *model.Filter) ([]*model.ItemOutput, error)
}

type Transaction interface {
	OpenShift(id int, transaction *model.Transaction) error
	CloseShift(id int, transaction *model.Transaction) error
	Collection(role string, id int, transaction *model.Transaction) error
	Income(role string, id int, transaction *model.Transaction) error
	PostavkaTransaction(role string, transaction *model.Transaction, postavkaID, shopID int) error
	GetAllShifts(filter *model.Filter) ([]*model.Shift, int64, error)
	GetShiftByID(id int) (*model.ShiftTransaction, error)
	GetLastShift(shopID int) (*model.Shift, error)
	GetAllTransaction(filter *model.Filter) ([]*model.TransactionResponse, int64, error)
	GetShiftByShopId(id int, workerId int) (*model.CurrentShift, error)
	GetTransactionByID(id int) (*model.TransactionResponse, error)
	UpdateTransaction(transaction *model.Transaction) error
	DeleteTransaction(id int) error
	CheckForBlockedShop(shiftID int) (bool, error)
}

type Master interface {
	GetShopsMaster() ([]*model.Shop, error)
	AddTovarMaster(tovar *model.TovarMaster) (*model.TovarMaster, error)
	NormaliseTovars() error
	NormaliseIngredients() error
	NormaliseTechCarts() error
	NormaliseNabors() error
	AddIngredientMaster(ingredient *model.IngredientMaster) (*model.IngredientMaster, error)
	AddTechCartMaster(techCart *model.TechCartMaster) (*model.TechCartMaster, error)
	AddOrUpdateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	AddNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	UpdateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error)
	DeleteNaborMaster(id int) error
	GetAllTovarMaster(filter *model.Filter) ([]*model.TovarMasterResponse, int64, error)
	GetTovarMaster(id int) (*model.TovarMasterResponse, error)
	GetAllIngredientMaster(filter *model.Filter) ([]*model.IngredientMasterResponse, int64, error)
	GetAllTechCartMaster(filter *model.Filter) ([]*model.TechCartMasterResponse, int64, error)
	GetIngredientMaster(id int) (*model.IngredientMasterResponse, error)
	GetTechCartMaster(id int) (*model.TechCartMasterResponse, error)
	GetAllNaborMaster(filter *model.Filter) ([]*model.NaborMasterOutput, int64, error)
	GetNaborMaster(id int) (*model.NaborMasterOutput, error)
	UpdateTovarMaster(tovar *model.ReqTovarMaster) error
	UpdateIngredientMaster(ingredient *model.ReqIngredientMaster) error
	UpdateTechCartMaster(techCart *model.ReqTechCartMaster) error
	DeleteTovarMaster(id int) error
	DeleteIngredientMaster(id int) error
	DeleteTechCartMaster(id int) error
	CreateNaborMaster(nabors []*model.NaborMaster) ([]*model.NaborMaster, error)
	UpdateTechCartsMaster(techCarts []*model.ReqTechCartMaster) error
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
	Registrate(req *model.ReqRegistrate) (string, error)
	VerifyEmail(req *model.ReqVerifyEmail) (int, string, error)
	SendCode(email string) error
	GetAllMobileUsers(filter *model.Filter) ([]*model.MobileUser, int64, error)
	GetMobileUser(id string, shopIDs []int) (*model.MobileUser, error)
	GetAllFeedbacks(filter *model.Filter) ([]*model.MobileUserFeedbackResponse, int64, error)
	GetFeedback(id int) (*model.MobileUserFeedbackResponse, error)
}

type Service struct {
	Tovar
	Terminal
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

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Terminal:    NewTerminalService(*repos),
		Tovar:       NewTovarService(repos.Tovar),
		Ingredient:  NewIngredientService(repos.Ingredient),
		Sklad:       NewSkladService(*repos),
		Dealer:      NewDealerService(repos.Dealer),
		Finance:     NewFinanceService(repos.Finance),
		Check:       NewCheckService(*repos),
		External:    NewExternalService(repos.External),
		User:        NewUserService(repos.User),
		Worker:      NewWorkerService(repos.Worker),
		Statistics:  NewStatisticsService(repos.Statistics),
		Transaction: NewTransactionService(*repos),
		Shop:        NewShopService(repos.Shop),
		Master:      NewMasterService(*repos),
		Mobile:      NewMobileService(repos.Mobile),
	}
}
