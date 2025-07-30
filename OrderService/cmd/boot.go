package cmd

import (
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"
	//"tesodev-korpes/pkg/middleware"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func BootOrderService(client *mongo.Client, e *echo.Echo) {
	config := config3.GetOrderConfig("dev")
	orderCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(orderCol)
	service := internal.NewService(repo)
	internal.NewHandler(e, service)

	e.Logger.Fatal(e.Start(config.Port))

}
