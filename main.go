package main

import (
	"github.com/labstack/echo/v4"
	"tesodev-korpes/CustomerService/cmd"
	"tesodev-korpes/pkg"
	"tesodev-korpes/shared/config"
)

func main() {
	//todo : what is dev,qa,prod ? explain why we are using them in the lecture
	dbConf := config.GetDBConfig("dev")

	client, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		panic(err)
	}

	e := echo.New()

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
