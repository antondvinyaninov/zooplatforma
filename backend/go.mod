module backend

go 1.25.5

require (
	database v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
	golang.org/x/crypto v0.46.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
)

replace database => ../database
