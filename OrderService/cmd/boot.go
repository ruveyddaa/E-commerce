// File: cmd/boot.go
package cmd

import (
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/client" // artık customerClient değil
	"time"

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

	// fasthttp tabanlı generic client (baseURL + timeout)
	cc := client.New("http://localhost:8001", 5*time.Second)

	// Service, HTTP çağrıları için generic client alıyor
	service := internal.NewService(repo, cc)

	internal.NewHandler(e, service)
	e.Logger.Fatal(e.Start(cfg.Port))
}
