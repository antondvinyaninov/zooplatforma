module clinic

go 1.25.5

require (
	database v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
)

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

replace database => ../../database
