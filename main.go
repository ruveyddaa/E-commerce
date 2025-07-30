package main

import (
	"fmt"
	"tesodev-korpes/CustomerService/cmd"
	ordercmd "tesodev-korpes/OrderService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/middleware"
	"tesodev-korpes/shared/config"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	//todo : what is dev,qa,prod ? explain why we are using them in the lecture
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

	orderEcho := echo.New()

	/*e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})*/

	customerEcho.GET("/swagger/*", echoSwagger.WrapHandler)
	// http://localhost:8001/swagger/index.html#/

	switch input {
	case "customer":
		cmd.BootCustomerService(client, customerEcho)
	case "order":

		ordercmd.BootOrderService(client, orderEcho)

	default:
		panic("Invalid input. Use 'customer', 'order', or 'both'.")
	}

	//challenge : after you create a func boot order service, manage somehow to run specific project
	//description : when you give an input here it should look that input and boot THAT specific project
	//if the input says "both" it should

	//PS : do not forget to create and call a different column for order service and do not forget to boot order service
	//from another port different than customer service

	//orderCol, err := pkg.GetMongoCollection(client, "tesodev", "order")
	//if err != nil {
	//	panic(err)

}
