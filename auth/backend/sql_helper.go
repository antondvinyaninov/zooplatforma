package main

import "fmt"

// sqlQuery возвращает SQL запрос с правильными placeholders для текущей БД
func sqlQuery(query string) string {
	if !isPostgres {
		return query
	}

	// Заменяем ? на $1, $2, $3 и т.д. для PostgreSQL
	result := ""
	paramNum := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result += fmt.Sprintf("$%d", paramNum)
			paramNum++
		} else {
			result += string(query[i])
		}
	}
	return result
}
