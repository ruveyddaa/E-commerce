package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
			fmt.Println(userRole)
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

	// --- DEBUG ADIMI: GELEN HAM JSON VERİSİNİ LOGLAMA ---
	// Cevabın body'sini okuyup bir byte dizisine alıyoruz.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cevap body'sini okurken hata: %w", err)
	}

	// Okuduğumuz ham veriyi string olarak konsola yazdırıyoruz.
	// SORUNUN KAYNAĞINI BURADA GÖRECEKSİNİZ!
	fmt.Println("MÜŞTERİ SERVİSİNDEN GELEN HAM JSON:", string(bodyBytes))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Müşteri servisi %d status döndürdü", resp.StatusCode)
	}

	var customer CustomerResponseModel
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&customer); err != nil {
		return nil, fmt.Errorf("JSON decode hatası: %w", err)
	}

	return &customer, nil
}

type CustomerResponseModel struct {
	ID        string    `bson:"_id" json:"id"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	Email     string    `bson:"email" json:"email"`
	Phone     []Phone   `bson:"phone" json:"phone"`
	Address   []Address `bson:"address" json:"address"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Role      string    `bson:"role" json:"role"` // role eklendi
}

type Address struct {
	Id      string `bson:"address_id,omitempty" json:"address_id"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
}

type Phone struct {
	Id          string `bson:"phone_id,omitempty" json:"phone_id"`
	PhoneNumber int    `bson:"phone_number" json:"phone_number"`
}
