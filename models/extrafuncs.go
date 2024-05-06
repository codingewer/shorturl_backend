package models

import (
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Rastgele isimoluşturan fonksiyon
func GenerateString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func ReplaceToStar(word string) string {
	lenght := len(word)
	staredPart := strings.Repeat("*", lenght-2)
	return word[:2] + staredPart
}
