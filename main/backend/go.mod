module backend

go 1.25.5

require (
	database v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/zooplatforma/pkg/clients v0.0.0-00010101000000-000000000000
)

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

replace database => ../../database

replace github.com/zooplatforma/pkg/clients => ../../pkg/clients
