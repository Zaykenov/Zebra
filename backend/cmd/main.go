package main

import (
	"errors"
	"log"
	"os"
	"time"
	"zebra"

	"zebra/model"
	"zebra/pkg/handler"
	"zebra/pkg/repository"
	"zebra/pkg/service"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Zebra Backend API
// @version 1.0
// @description API server for zebra-crm system
// @host localhost:4000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Print("No .env file found, please set")
	}
}

func main() {
	logrus.Print("Startup server")
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initEnv(); err != nil {
		logrus.Fatalf("error initializing env: %s", err.Error())
	}

	db, err := repository.NewPostgreDB(
		os.Getenv("DSN"),
	)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = migrateGorm(gormDB)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	repos := repository.NewRepository(db, gormDB)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)
	staticHandler := handler.NewStaticHandler(service)
	srv := new(zebra.Server)
	staticSrv := new(zebra.Server)

	logrus.Print("Server Runing on Dev mode")

	go sendToTis(service)
	go DailyStatistic(service)
	go RecalculateTrafficReport(service)

	go func() {
		if err := srv.Run(os.Getenv("APIPortHTTP"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	if err := staticSrv.Run(os.Getenv("StaticPortHTTP"), staticHandler.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func migrateGorm(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Shop{},
		&model.Tovar{},
		&model.TechCartMaster{},
		&model.Sklad{},
		&model.SkladTovar{},
		&model.CategoryTovar{},
		&model.Ingredient{},
		&model.SkladIngredient{},
		&model.CategoryIngredient{},
		&model.Schet{},
		&model.Dealer{},
		&model.User{},
		&model.UserQR{},
		&model.Feedback{},
		&model.Check{},
		&model.NaborMaster{},
		&model.Nabor{},
		&model.TechCart{},
		&model.CheckTovar{},
		&model.CheckTechCart{},
		&model.CheckModificator{},
		&model.CheckView{},
		&model.Postavka{},
		&model.IngredientTechCart{},
		&model.NaborTechCart{},
		&model.IngredientNabor{},
		&model.ItemPostavka{},
		&model.RemoveFromSklad{},
		&model.RemoveFromSkladItem{},
		&model.Tag{},
		&model.Worker{},
		&model.TisResponse{},
		&model.Shift{},
		&model.Transaction{},
		&model.TransactionPostavka{},
		&model.Transfer{},
		&model.ItemTransfer{},
		&model.Inventarization{},
		&model.InventarizationItem{},
		&model.ExpenceTovar{},
		&model.ExpenceIngredient{},
		&model.ErrorCheck{},
		&model.SendToTis{},
		&model.DailyStatistic{},
		&model.InventarizationGroup{},
		&model.InventarizationGroupItem{},
		&model.AsyncJob{},
		&model.TovarMaster{},
		&model.IngredientMaster{},
		&model.IngredientTechCartMaster{},
		&model.IngredientNaborMaster{},
		&model.NaborTechCartMaster{},
		&model.Client{},
		&model.FailedCheck{},
		&model.Stolik{},
	)
	if err != nil {
		return err
	}

	return nil
}

func initEnv() error {
	reqs := []string{}
	for i := 0; i < len(reqs); i++ {
		_, exists := os.LookupEnv(reqs[i])
		if !exists {
			return errors.New(reqs[i] + ".env variables not set")
		}
	}
	return nil
}

func sendToTis(service *service.Service) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minute().Do(func() {
		log.Print("Sending to tis")
		err := service.SendToTis()
		if err != nil {
			logrus.Print(err)
		}
	})

	s.StartAsync()
}

func RecalculateTrafficReport(service *service.Service) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(30).Minute().Do(func() {
		err := service.RecalculateTrafficReport()
		if err != nil {
			logrus.Print(err)
			service.SaveError(err.Error(), "RecalculateTrafficReport")
		}
	})

	s.StartAsync()
}

func DailyStatistic(service *service.Service) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At("18:00:00").Do(func() {
		err := service.DailyStatistic()
		if err != nil {
			logrus.Print(err)
		}
	})

	s.StartAsync()
}
