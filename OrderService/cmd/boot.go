package cmd

import (
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func BootOrderService(clientMongo *mongo.Client, e *echo.Echo) {
	cfg := config3.GetOrderConfig("dev")

	orderCol, err := pkg.GetMongoCollection(clientMongo, cfg.DbConfig.DBName, cfg.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(orderCol)

	// Sadece base URL veriyoruz (service.go net/http versiyonu bunu bekliyor)
	customerBaseURL := "http://localhost:8001"
	service := internal.NewService(repo, customerBaseURL)

	internal.NewHandler(e, service)
	e.Logger.Fatal(e.Start(cfg.Port))
}
