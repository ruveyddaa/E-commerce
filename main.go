package main

import (
	"fmt"
	"tesodev-korpes/CustomerService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/middleware"
	"tesodev-korpes/shared/config"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func main() {
	//todo : what is dev,qa,prod ? explain why we are using them in the lecture
	dbConf := config.GetDBConfig("dev")

	client, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("connecting db")

	e := echo.New()

	e.Use(middleware.CorrelationIdMiddleware())
	e.Use(middleware.LoggingMiddleware)
	e.Use(middleware.RecoveryMiddleware)
	e.Use(middleware.ErrorHandler())

	/*e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})*/

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// http://localhost:8001/swagger/index.html#/

	cmd.BootCustomerService(client, e)

	//challenge : after you create a func boot order service, manage somehow to run specific project
	//description : when you give an input here it should look that input and boot THAT specific project
	//if the input says "both" it should

	//PS : do not forget to create and call a different column for order service and do not forget to boot order service
	//from another port different than customer service

	//orderCol, err := pkg.GetMongoCollection(client, "tesodev", "order")
	//if err != nil {
	//	panic(err)

}
