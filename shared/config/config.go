package config

import (
	"time"
)

// explain why we have the "shared" folder, why we have a config here and another config in seperate projects in the lecture?
type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://root:root1234@mongodb_docker:27017",
	},
	"qa": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://root:root1234@mongodb_docker:27017",
	},
	"dev": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://localhost:27017/",
	},
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
