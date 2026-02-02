#!/bin/bash
# Скрипт для исправления SQL запросов в handlers.go

FILE="handlers.go"

# Функция для замены SQL запросов
fix_sql() {
    # Заменяем все WHERE ... = ? на условные блоки
    perl -i -pe '
        # Заменяем простые SELECT с одним параметром
        s/db\.QueryRow\(`\s*SELECT (.*?) FROM (\w+) WHERE (\w+) = \?\s*`, ([^)]+)\)/db.QueryRow(func() string { if isPostgres { return "SELECT $1 FROM $2 WHERE $3 = \$1" } else { return "SELECT $1 FROM $2 WHERE $3 = ?" } }(), $4)/g;
        
        # Заменяем Query с одним параметром
        s/db\.Query\("SELECT (.*?) FROM (\w+) WHERE (\w+) = \?", ([^)]+)\)/db.Query(func() string { if isPostgres { return "SELECT $1 FROM $2 WHERE $3 = \$1" } else { return "SELECT $1 FROM $2 WHERE $3 = ?" } }(), $4)/g;
    ' "$FILE"
}

fix_sql
echo "SQL queries fixed"
