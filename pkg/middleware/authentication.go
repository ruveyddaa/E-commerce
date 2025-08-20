package middleware

import (
	"context"
	"fmt"
	"strings"
	"tesodev-korpes/pkg/auth"
	"tesodev-korpes/pkg/customError"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

			role, err := ExistUser(mongoClient, claims.ID)
			if err != nil {
				return customError.NewUnauthorized(customError.MissingAuthToken)
			}
			fmt.Println("fjksdjkfjksdfkjsdfjkjndsfjfdsjkjkfds", claims.ID, role.Membership, role.SystemRole)
			c.Set("userId", claims.ID)
			c.Set("userRole", role.SystemRole)
			c.Set("userMembership", role.Membership)
			return next(c)
		}
	}
}

type RoleInfo struct {
	SystemRole string `bson:"role"`
	Membership string `bson:"membership"`
}

func ExistUser(mongoClient *mongo.Client, userID string) (*RoleInfo, error) {
	col := mongoClient.Database("tesodev").Collection("customer")

	// 1. Decode edilecek struct, veritabanından gelecek olan yapıyla eşleşmeli.
	// Projeksiyon kullandığımız için bize sadece { "role": { ... } } yapısında bir döküman gelecek.
	var result struct {
		Role RoleInfo `bson:"role"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Sadece "role" alanını çekmek için projeksiyon (projection) tanımla.
	// Bu, gereksiz veri transferini önler ve performansı artırır.
	opts := options.FindOne().SetProjection(bson.M{"role": 1})

	// 3. Sorguya projeksiyonu ekleyerek FindOne işlemini yap ve sonucu "result" struct'ına decode et.
	err := col.FindOne(ctx, bson.M{"_id": userID}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, customError.NewNotFound(customError.CustomerNotFound)
		}
		return nil, err
	}

	// 4. Dönen sonucun içindeki gömülü Role objesini (artık RoleInfo tipinde) döndür.
	fmt.Printf("Role found: Membership=%s, SystemRole=%s\n", result.Role.Membership, result.Role.SystemRole)
	return &result.Role, nil
}
