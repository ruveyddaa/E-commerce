package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func AuthorizationMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("userID").(string)
			if !ok || userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Kullanıcı doğrulanamadı")
			}
			userRole, ok := c.Get("userRole").(string)
			if !ok || userRole == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Role alınımadı")
			}

			for _, allowed := range allowedRoles {
				if strings.EqualFold(userRole, allowed) {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Bu işlemi yapmaya yetkiniz yok")
		}
	}
}

func fetchCustomerByID(customerServiceURL, customerID string) (*CustomerResponseModel, error) {
	url := fmt.Sprintf("%s/customer/%s", customerServiceURL, customerID)
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Müşteri servisi %d status döndürdü", resp.StatusCode)
	}

	var customer CustomerResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

type CustomerResponseModel struct {
	ID        string `bson:"_id" json:"id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
	IsActive  bool   `bson:"is_active" json:"is_active"`
	Role      string `bson:"role" json:"role"` // role eklendi
}
