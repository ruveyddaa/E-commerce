package authentication

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

/*
var SecretKey = []byte("secret-key")

	type Claims struct {
		Id        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		jwt.StandardClaims
	}

	func CreateJWT(Id string, firstName string, lastName string, key string) string {
		claims := &Claims{
			Id:        Id,
			FirstName: firstName,
			LastName:  lastName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(key))
		if err != nil {
			return err.Error()
		}
		return tokenString
	}

	func VerifyJWT(tokenString string) (*Claims, error) {
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			return nil, echo.ErrUnauthorized
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, echo.ErrUnauthorized
		}
		return claims, nil
	}
*/
func HashPassword(password []byte) ([]byte, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
func CheckPasswordHash(plainPassword string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
