package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"tesodev-korpes/pkg/auth"
	"tesodev-korpes/pkg/customError"
	"time"
)

type SkipperFunc func(c echo.Context) bool

func Authentication(mongoClient *mongo.Client, skipper SkipperFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper != nil && skipper(c) {
				return next(c)
			}
			const bearerPrefix = "Bearer "

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
				return customError.NewUnauthorized(customError.MissingAuthToken)
			}

			tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := auth.VerifyJWT(tokenStr)
			if err != nil {
				return customError.NewUnauthorized(customError.MissingAuthToken)
			}

			c.Set("userID", claims.ID)
			role, err := getUserRoleFromMongo(mongoClient, claims.ID)
			if err != nil {
				return customError.NewUnauthorized("")
			}
			c.Set("role", role)
			return next(c)
		}
	}
}

func getUserRoleFromMongo(mongoClient *mongo.Client, userID string) (string, error) {
	col := mongoClient.Database("tesodev").Collection("customer")

	var doc struct {
		System struct {
			Role string `bson:"role"`
		} `bson:"system"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := col.FindOne(ctx, bson.M{"_id": userID}).Decode(&doc); err != nil {
		return "", err
	}

	if doc.System.Role == "" {
		return "", fmt.Errorf("role not found for user %s", userID)
	}

	return doc.System.Role, nil
}
