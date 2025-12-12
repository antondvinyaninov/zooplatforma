package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := "jyjy4VlgOPGIPSG5vJPurXDnd1ZpHj2X2dIBtdWfjJE="

	fmt.Println("Секрет:", secret)
	fmt.Println("Длина:", len(secret))
	fmt.Println("")

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"email":   "anton@dvinyaninov.ru",
		"roles":   []string{"user", "superadmin"},
		"exp":     1766149090,
		"iat":     1765544290,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	fmt.Println("Созданный токен:")
	fmt.Println(tokenString)
	fmt.Println("")

	// Токен из curl
	curlToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFudG9uQGR2aW55YW5pbm92LnJ1IiwiZXhwIjoxNzY2MTQ5MDkwLCJpYXQiOjE3NjU1NDQyOTAsInJvbGVzIjpbInVzZXIiLCJzdXBlcmFkbWluIl0sInVzZXJfaWQiOjF9.aTzkXOIvtdTbNIXBTSItKDVQtiFA2Adv-AZUVLlDE7E"

	fmt.Println("Токен из curl:")
	fmt.Println(curlToken)
	fmt.Println("")

	// Парсим токен из curl
	parsed, err := jwt.Parse(curlToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("❌ Ошибка парсинга:", err)

		// Попробуем без = в конце
		secretWithout := secret[:len(secret)-1]
		fmt.Println("")
		fmt.Println("Пробуем без '=' в конце:", secretWithout)

		parsed2, err2 := jwt.Parse(curlToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretWithout), nil
		})

		if err2 != nil {
			fmt.Println("❌ Тоже ошибка:", err2)
		} else if parsed2.Valid {
			fmt.Println("✅ Работает без '='!")
		}

		os.Exit(1)
	}

	if parsed.Valid {
		fmt.Println("✅ Токен валиден!")
		claims := parsed.Claims.(jwt.MapClaims)
		fmt.Println("Claims:", claims)
	} else {
		fmt.Println("❌ Токен невалиден!")
	}
}
