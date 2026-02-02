package handlers

import (
	"fmt"
	"os"
)

// ConvertPlaceholders converts ? to $1, $2, $3 for PostgreSQL in production
func ConvertPlaceholders(query string) string {
	if os.Getenv("ENVIRONMENT") == "production" {
		result := ""
		paramNum := 1
		for _, char := range query {
			if char == '?' {
				result += fmt.Sprintf("$%d", paramNum)
				paramNum++
			} else {
				result += string(char)
			}
		}
		return result
	}
	return query
}
