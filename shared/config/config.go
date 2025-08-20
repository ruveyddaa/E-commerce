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

type AuthConfig struct {
	JWTSecret string
}

type Config struct {
	RoleMapping       map[string]string
	EndpointRoles     map[string][]string
	EndpointRolesPath map[string][]string
}

var Cfg = Config{
	EndpointRoles: map[string][]string{
		"CustomerRead":   {"admin", "manager", "user"},
		"CustomerWrite":  {"admin", "manager"},
		"CustomerDelete": {"admin"},
		"CustomerList":   {"admin", "manager"},

		"OrderCreate": {"premium", "non-premium", "user"},
		"OrderRead":   {"admin", "manager", "user"},
		"OrderUpdate": {"admin", "manager"},
		"OrderCancel": {"admin", "manager", "user"},
		"OrderList":   {"admin", "manager"},
	},
	EndpointRolesPath: map[string][]string{

		"POST /customer/create": []string{},
		"POST /customer/login":  []string{},

		"GET /customer/verify":       {"admin"},
		"GET /customer/list":         {"admin", "manager", "user"},
		"GET /customer/:id":          {"admin", "manager", "user"},
		"GET /customer/email/:email": {"admin", "manager", "user"},
		"PUT /customer/:id":          {"admin", "manager", "user"},
		"DELETE /customer/:id":       {"admin"},

		"POST /order":              {"admin", "manager", "user"},
		"GET /order/list":          {"admin", "manager", "user"},
		"GET /order/:id":           {"admin", "manager", "user"},
		"PATCH /order/:id/ship":    {"admin", "manager", "user"},
		"PATCH /order/:id/deliver": {"admin", "manager", "user"},
		"DELETE /order/cancel/:id": {"admin", "manager", "user"},
	},
	RoleMapping: map[string]string{
		"premium":     "/internal/price/premium/:id",
		"non-premium": "/internal/price/non-premium/:id",
	},
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

var (
	serviceUrls ServiceUrls
	authConfig  AuthConfig
)

func init() {

	if err := godotenv.Load("./media/.env"); err != nil {
		panic("Environment variable did not load")
	}
	fmt.Println("Environment variables loaded")

	serviceUrls = ServiceUrls{
		CustomerServiceURL: os.Getenv("CUSTOMER_SERVICE_URL"),
	}
	authConfig = AuthConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if authConfig.JWTSecret == "" {
		panic("JWT_SECRET environment variable not set")
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

func GetAuthConfig() AuthConfig {
	return authConfig
}
