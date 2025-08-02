package authentication

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
func CheckPasswordHash(plainPassword []byte, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
