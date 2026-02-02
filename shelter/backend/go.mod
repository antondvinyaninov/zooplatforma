module shelter

go 1.25.5

replace database => ../../database

replace github.com/zooplatforma/pkg/middleware => ../../pkg/middleware

require (
	database v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/joho/godotenv v1.5.1
	github.com/zooplatforma/pkg/middleware v0.0.0-00010101000000-000000000000
)

require (
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.33 // indirect
)
