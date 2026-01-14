module admin

go 1.25.5

replace database => ../../database

require (
	database v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
)

require github.com/mattn/go-sqlite3 v1.14.32 // indirect
