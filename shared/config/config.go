package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// explain why we have the "shared" folder, why we have a config here and another config in seperate projects in the lecture?
type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: MongoUrlLoad(),
	},
	"qa": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: MongoUrlLoad(),
	},
	"dev": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: MongoUrlLoad(),
	},
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}


	return &config
}

func MongoUrlLoad() string {
	if err := godotenv.Load("./media/.env"); err != nil {
		panic("Environment variable did not load")
	}
	fmt.Println("Connected environment variable load")

	return os.Getenv("MONGO_URI")
}
