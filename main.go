package main

import (
	"fmt"
	customercmd "tesodev-korpes/CustomerService/cmd"
	"tesodev-korpes/CustomerService/controller"
	ordercmd "tesodev-korpes/OrderService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/middleware"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	dbConf := config.GetDBConfig("dev")

	client, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("connecting db")

	customerEcho := echo.New()
	customerEcho.Use(middleware.CorrelationIdMiddleware())
	customerEcho.Use(middleware.LoggingMiddleware)
	customerEcho.Use(middleware.RecoveryMiddleware)
	customerEcho.Use(middleware.ErrorHandler())
	customerEcho.GET("/swagger/*", echoSwagger.WrapHandler)
	// customerEcho üzerinde /price endpoint’i
	customerEcho.GET("/price",
		func(c echo.Context) error { return nil },                   // boş handler, middleware yönlendirecek
		middleware.Authentication(client, nil),                      // senin hazır authentication
		middleware.AuthorizationMiddleware(config.Cfg.AllowedRoles), // senin authorization
		middleware.RoleRouting(config.Cfg),                          // role routing middleware
	)

	// internal handler'lar (dışarıdan görünmez)
	customerEcho.GET("/internal/price/premium", controller.HandlePremiumPrice)
	customerEcho.GET("/internal/price/non-premium", controller.HandleNonPremiumPrice)

	orderEcho := echo.New()
	orderEcho.Use(middleware.CorrelationIdMiddleware())
	orderEcho.Use(middleware.LoggingMiddleware)
	orderEcho.Use(middleware.RecoveryMiddleware)
	orderEcho.Use(middleware.ErrorHandler())
	orderEcho.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		customercmd.BootCustomerService(client, customerEcho)
	}()

	ordercmd.BootOrderService(client, orderEcho)
}
