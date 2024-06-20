package handler

import (
	"zebra/pkg/service"

	_ "zebra/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/signin", h.signIn)
	router.POST("/upload", h.uploadImage)
	router.GET("/itemImage/:fileName", h.ItemImageHandler)
	router.GET("check/getForPrinter/:id", h.getCheckByIDForPrinter)
	router.GET("/ping", h.PingPong)
	worker := router.Group("", h.workerIdentity)
	{
		worker.POST("/terminal", h.Terminal)
		shop := worker.Group("/shop")
		{
			shop.GET("/getAll", h.getAllShops)
			shop.GET("/get/:id", h.getShopByID)
		}
		finance := worker.Group("/finance")
		{
			schet := finance.Group("/schet")
			{
				schet.GET("/getAll", h.getAllSchet)
			}
		}
		transaction := worker.Group("/transaction")
		{
			transaction.POST("/create", h.createTransaction)
			transaction.GET("/getAll", h.getAllTransaction)
			transaction.GET("/get/:id", h.getTransactionByID)
			transaction.POST("/update", h.updateTransaction)
			transaction.POST("/delete/:id", h.deleteTransaction)
		}
		shifts := worker.Group("/shift")
		{
			shifts.GET("/getAll", h.getAllShifts)
			shifts.GET("/getAll/:page", h.getAllShifts)
			shifts.GET("/get/:id", h.getShiftByID)
			shifts.GET("/check", h.getShiftByShopId)
		}
		terminal := worker.Group("/terminal")
		{
			terminal.GET("/start", h.startTerminal)
		}
		tovar := worker.Group("/item")
		{
			tovar.GET("/getAll", h.getAllTovar)
			tovar.GET("/get", h.getTovar)
			tovar.GET("/getWithParams", h.getTovarWithParams)
			tovar.GET("/getItems", h.getItems)
			tovar.GET("/getEverything", h.getEverything)
			tovar.POST("/recalculate", h.RecalculateNetCost)
			categoryTovar := tovar.Group("/category")
			{
				categoryTovar.GET("/getAll", h.getAllCategoryTovar)
				categoryTovar.GET("/get", h.getCategoryTovar)
			}
			techCart := tovar.Group("/techCart")
			{
				techCart.GET("/getAll", h.getAllTechCart)
				techCart.GET("/get", h.getTechCart)
				techCart.GET("/getNabor", h.getTechCartNabor)
				techCart.GET("/getWithParams", h.getTechCartWithParams)
			}
		}
		ingredient := worker.Group("/ingredient")
		{
			ingredient.GET("/getAll", h.getAllIngredient)
			ingredient.GET("/get", h.getIngredient)
			categoryIngredient := ingredient.Group("/category")
			{
				categoryIngredient.GET("/getAll", h.getAllCategoryIngredient)
				categoryIngredient.GET("/get", h.getCategoryIngredient)
			}
			nabor := ingredient.Group("/nabor")
			{
				nabor.GET("/getAll", h.getAllNabor)
				nabor.GET("/get", h.getNabor)
			}
		}
		sklad := worker.Group("/sklad")
		{
			sklad.GET("/getAll", h.getAllSklad)
			sklad.POST("/remove", h.RemoveFromSklad)
			postavka := sklad.Group("/postavka")
			{
				postavka.POST("/create", h.addToSklad)
				postavka.GET("/getAll/:page", h.GetAllPostavka)
				postavka.GET("/getAll/", h.GetAllPostavka)
				postavka.GET("/get/:id", h.GetPostavka)
				postavka.POST("/update", h.UpdatePostavka)
				postavka.POST("/delete", h.DeletePostavka)
				postavka.POST("/createWorker", h.addToSkladWorker)
			}
			inventratization := sklad.Group("/inventarization")
			{
				inventratization.GET("/getAll", h.GetAllInventarization)
				inventratization.POST("/create", h.CreateInventarization)
				inventratization.GET("/get/:id", h.GetInventarization)
				inventratization.GET("/recalculate", h.RecalculateInventarization)
				inventratization.POST("/update", h.UpdateInventarization)
				inventratization.POST("/updateParams", h.UpdateInventarizationParams)

				group := inventratization.Group("/group")
				{
					group.GET("/getAll", h.GetAllInventarizationGroup)
					group.POST("/create", h.CreateInventarizationGroup)
					group.GET("/get/:id", h.GetInventarizationGroup)
					group.GET("/deleteItem/:id", h.DeleteInventarizationGroupItem)
					group.POST("/update", h.UpdateInventarizationGroup)
				}
			}
			transfer := sklad.Group("/transfer")
			{
				transfer.POST("/create", h.AddTransfer)
				transfer.GET("/getAll/", h.GetAllTransfer)
				transfer.GET("/get/:id", h.GetTransfer)
				transfer.POST("/update", h.UpdateTransfer)
				transfer.POST("/delete", h.DeleteTransfer)
			}
		}
		dealer := worker.Group("/dealer")
		{
			dealer.GET("/getAll", h.getAllDealer)
		}
		check := worker.Group("/check")
		{
			check.GET("/getStoliki", h.getStoliki)
			check.GET("/getFilledStoliki", h.getFilledStoliki)
			check.POST("/create", h.addCheck)
			check.POST("/failed", h.saveFailedCheck)
			check.GET("/getAll/:page", h.getAllCheck)
			check.GET("/getAll", h.getAllCheck)
			check.GET("/getAllWorker", h.getAllCheckWorker)
			check.GET("/get/:id", h.getCheckByID)
			check.POST("/delete", h.DeleteCheck)
			check.POST("/idempotencyCheck", h.IdempotencyCheck)
			tag := check.Group("/tag")
			{
				tag.POST("/create", h.addTag)
				tag.GET("/getAll", h.getAllTag)
				tag.GET("/get/:id", h.getTag)
				tag.POST("/update", h.updateTag)
				tag.POST("/delete", h.deleteTag)
			}
		}
		external := worker.Group("/external")
		{
			external.POST("/saveCheck", h.saveCheck)
		}

	}
	router.POST("/authorize", h.authorize)

	admin := router.Group("", h.adminIdentity)
	{
		mobile := admin.Group("/mobile")
		{
			user := mobile.Group("/user")
			{
				user.GET("/getAll", h.getAllMobileUsers)
				user.GET("/get/:id", h.getMobileUser)
			}
			feedback := mobile.Group("/feedback")
			{
				feedback.GET("/getAll", h.getAllFeedbacks)
				feedback.GET("/get/:id", h.getFeedback)
			}
		}
		master := admin.Group("/master")
		{
			master.GET("/item/getAll", h.getAllTovarMaster)
			master.GET("/item/get/:id", h.getTovarMaster)
			master.GET("/ingredient/getAll", h.getAllIngredientMaster)
			master.GET("/ingredient/get/:id", h.getIngredientMaster)
			master.GET("/techCart/getAll", h.getAllTechCartMaster)
			master.GET("/techCart/get/:id", h.getTechCartMaster)
			master.GET("/techCart/updateAll", h.updateTechCartsMaster)
			master.GET("/ingredient/nabor/getAll", h.getAllNaborMaster)
			master.GET("/ingredient/nabor/get/:id", h.getNaborMaster)
			master.POST("/ingredient/nabor/createAll", h.createNaborMaster)
		}
		shop := admin.Group("/shop")
		{
			shop.POST("/create", h.createShop)
			shop.POST("/update", h.updateShop)
			shop.POST("/delete", h.deleteShop)
		}
		tovar := admin.Group("/item")
		{
			tovar.POST("/create", h.addTovar)

			tovar.POST("/update", h.updateTovar)
			tovar.POST("/delete", h.deleteTovar)

			categoryTovar := tovar.Group("/category")
			{
				categoryTovar.POST("/delete", h.deleteCategoryTovar)
			}
			techCart := tovar.Group("/techCart")
			{
				techCart.POST("/create", h.addTechCart)
				techCart.POST("/update", h.updateTechCart)
				techCart.POST("/delete", h.deleteTechCart)

			}
		}
		ingredient := admin.Group("/ingredient")
		{
			ingredient.POST("/update", h.updateIngredient)
			ingredient.POST("/create", h.addIngredient)
			ingredient.POST("/getTechCartByIngredientID", h.getTechCartByIngredientID)
			ingredient.POST("/delete", h.deleteIngredient)
			categoryIngredient := ingredient.Group("/category")
			{
				categoryIngredient.POST("/delete", h.deleteCategoryIngredient)
			}
			nabor := ingredient.Group("/nabor")
			{
				nabor.POST("/createAll", h.addNabors)
				nabor.POST("/delete", h.deleteNabor)
			}
		}
		sklad := admin.Group("/sklad")
		{
			sklad.POST("/create", h.addSklad)
			sklad.GET("/get", h.getSklad)
			sklad.POST("/update", h.updateSklad)
			sklad.POST("/delete", h.deleteSklad)
			sklad.GET("/ostatki", h.ostatki)
			sklad.GET("/getRemoved", h.GetRemoved)
			sklad.GET("/getRemoved/:id", h.GetRemovedByID)
			sklad.GET("/getTrafficReport", h.GetTrafficReport)
			sklad.POST("/recalculateDailyStatistic", h.RecalculateDailyStatistic)
			sklad.POST("/confirm/:id", h.ConfirmToRemove)
			sklad.POST("/reject/:id", h.RejectToRemove)
			inventratization := sklad.Group("/inventarization")
			{
				inventratization.POST("/delete", h.DeleteInventarization)
				inventratization.POST("/recalculate/:id", h.RecalculateOneInventarization)
				inventratization.GET("/getDetails/income/:id", h.GetInventarizationDetailsIncome)
				inventratization.GET("/getDetails/expence/:id", h.GetInventarizationDetailsExpence)
				inventratization.GET("/getDetails/spisanie/:id", h.GetInventarizationDetailsSpisanie)
				inventratization.POST("/deleteItem", h.DeleteInventarizationItem)
			}
			spisanie := sklad.Group("/spisanie")
			{
				spisanie.POST("/update", h.UpdateSpisanie)
				spisanie.POST("/delete", h.DeleteSpisanie)
			}
		}
		dealer := admin.Group("/dealer")
		{
			dealer.POST("/create", h.addDealer)
			dealer.GET("/get", h.getDealer)
			dealer.POST("/update", h.updateDealer)
			dealer.POST("/delete", h.deleteDealer)
		}
		check := admin.Group("/check")
		{
			check.GET("/deactivate/:id", h.DeactivateCheck)
		}
		finance := admin.Group("/finance")
		{
			schet := finance.Group("/schet")
			{
				schet.POST("/create", h.addSchet)
				schet.GET("/get", h.getSchet)
				schet.POST("/update", h.updateSchet)
				schet.POST("/delete", h.deleteSchet)
			}
		}
		statistics := admin.Group("/statistics")
		{
			statistics.GET("/workers", h.getWorkersStat)
			statistics.GET("/today", h.todayStatistics)
			statistics.GET("/everyDay", h.everyDayStatistics)
			statistics.GET("/everyWeek", h.everyWeekStatistics)
			statistics.GET("/everyMonth", h.everyMonthStatistics)
			statistics.GET("/payments", h.payments)
			statistics.GET("/daysOfTheWeek", h.DaysOfTheWeek)
			statistics.GET("/statByHour", h.StatByHour)
			statistics.GET("/ABC", h.ABC)
			statistics.GET("/topSales", h.TopSales)

		}
		workers := admin.Group("/workers")
		{
			workers.GET("/getAll", h.getAllWorkers)
			workers.GET("/get/:id", h.getWorker)
			workers.POST("/update", h.updateWorker)
			workers.POST("/delete/:id", h.deleteWorker)
		}
	}
	master := router.Group("/master", h.masterIdentity)
	{
		master.GET("/getShops", h.GetShopsMaster)
		master.GET("/normaliseTovars", h.NormaliseTovars)
		master.GET("/normaliseIngredients", h.NormaliseIngredients)
		master.GET("/normaliseTechCarts", h.NormaliseTechCarts)
		master.GET("/normaliseNabors", h.NormaliseNabors)
		master.POST("/item/update", h.UpdateTovarMaster)
		master.POST("/item/confirm/:id", h.ConfirmTovarMaster)
		master.POST("/item/reject/:id", h.RejectTovarMaster)
		master.POST("/item/delete", h.DeleteTovarMaster)
		master.POST("/item/techCart/create", h.addTechCart)
		master.POST("/techCart/update", h.UpdateTechCartMaster)
		master.POST("/techCart/confirm/:id", h.ConfirmTechCartMaster)
		master.POST("/techCart/reject/:id", h.RejectTechCartMaster)
		master.POST("/techCart/delete", h.DeleteTechCartMaster)
		master.POST("/ingredient/update", h.UpdateIngredientMaster)
		master.POST("/ingredient/confirm/:id", h.ConfirmIngredientMaster)
		master.POST("/ingredient/reject/:id", h.RejectIngredientMaster)
		master.POST("/ingredient/delete", h.DeleteIngredientMaster)
		master.POST("/ingredient/nabor/add", h.addNaborMaster)
		master.POST("/ingredient/nabor/update", h.updateNaborMaster)
		master.POST("/ingredient/nabor/confirm/:id", h.ConfirmNaborMaster)
		master.POST("/ingredient/nabor/reject/:id", h.RejectNaborMaster)
		master.POST("/ingredient/nabor/delete", h.DeleteNaborMaster)
		master.POST("/item/category/create", h.addCategoryTovar)
		master.POST("/item/category/update", h.updateCategoryTovar)
		master.POST("/ingredient/category/create", h.addCategoryIngredient)
		master.POST("/ingredient/category/update", h.updateCategoryIngredient)
	}

	user := router.Group("/user")
	{
		user.POST("/create", h.addUser)
		user.GET("/getInfo/:id", h.getUserInfo)
		user.GET("/getCurrentOrders/:id", h.getCurrentOrders)
		user.POST("/addFeedback", h.addFeedback)
		user.GET("/genCode/:id", h.genCode)
		user.GET("/getUserByCode/:code", h.getUserByCode)
	}

	router.POST("/registrate", h.Registrate)
	router.POST("/registrate/verify-email", h.VerifyEmail)
	router.POST("/sign-in", h.signInClient)
	router.POST("/sign-in/verify-email", h.VerifyEmail)

	return router
}
