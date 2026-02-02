#!/bin/bash

# Скрипт для инициализации Auth Service

echo "🚀 Initializing Auth Service..."
echo ""

# Перейти в директорию auth/backend
cd "$(dirname "$0")/../auth/backend" || exit 1

# Создать директорию для базы данных
echo "📁 Creating database directory..."
mkdir -p database

# Создать файл базы данных
echo "💾 Creating database file..."
touch database/auth.db

# Установить Go зависимости
echo "📦 Installing Go dependencies..."
go mod download

echo ""
echo "✅ Auth Service initialized successfully!"
echo ""
echo "To start the service, run:"
echo "  cd auth/backend"
echo "  air"
echo ""
echo "Or use:"
echo "  make all"
