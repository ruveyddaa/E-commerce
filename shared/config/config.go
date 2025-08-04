package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
}

type ServiceUrls struct {
	CustomerServiceURL string
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration: time.Second * 100,
	},
	"qa": {
		MongoDuration: time.Second * 100,
	},
	"dev": {
		MongoDuration: time.Second * 100,
	},
}

var serviceUrls ServiceUrls

func init() {

	if err := godotenv.Load("./media/.env"); err != nil {
		panic("Environment variable did not load")
	}
	fmt.Println("Environment variables loaded")

	serviceUrls = ServiceUrls{
		CustomerServiceURL: os.Getenv("CUSTOMER_SERVICE_URL"),
	}
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}

	config.MongoClientURI = os.Getenv("MONGO_URI")
	return &config
}

func GetServiceURLs() ServiceUrls {
	return serviceUrls
}
