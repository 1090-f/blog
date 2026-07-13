package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
