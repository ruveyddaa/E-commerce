package config

type CustomerConfig struct {
	Port     string
	DbConfig struct {
		DBName  string
		ColName string
	}
}

var cfgs = map[string]CustomerConfig{
	"prod": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
	"qa": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
	"dev": {
		Port: ":8001",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "customer",
		},
	},
}

func GetCustomerConfig(env string) *CustomerConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
