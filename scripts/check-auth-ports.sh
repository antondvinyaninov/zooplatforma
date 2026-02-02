#!/bin/bash

# Скрипт для проверки правильности использования портов для авторизации
# Auth Service должен быть на порту 7100, НЕ 8000!

echo "🔍 Проверка использования портов для авторизации..."
echo ""

ERRORS=0

# Проверяем frontend файлы на использование 8000 для auth
echo "📁 Проверка main/frontend..."
if grep -r "localhost:8000.*auth" main/frontend --include="*.ts" --include="*.tsx" 2>/dev/null; then
    echo "❌ ОШИБКА: Найдено использование localhost:8000 для авторизации!"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ Нет прямых ссылок на localhost:8000 для auth"
fi

# Проверяем использование apiUrl для auth
echo ""
echo "📁 Проверка использования apiUrl для auth..."
if grep -r "apiUrl.*auth" main/frontend --include="*.ts" --include="*.tsx" 2>/dev/null; then
    echo "⚠️  ВНИМАНИЕ: Найдено использование apiUrl для auth (должен быть authUrl)"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ Нет использования apiUrl для auth"
fi

# Проверяем что authApi использует authClient
echo ""
echo "📁 Проверка authApi в lib/api.ts..."
if grep -A 5 "export const authApi" main/frontend/lib/api.ts | grep -q "authClient"; then
    echo "✅ authApi использует authClient (порт 7100)"
else
    echo "❌ ОШИБКА: authApi НЕ использует authClient!"
    ERRORS=$((ERRORS + 1))
fi

# Проверяем что authClient указывает на 7100
echo ""
echo "📁 Проверка authClient URL..."
if grep "AUTH_URL.*7100" main/frontend/lib/api.ts >/dev/null 2>&1; then
    echo "✅ AUTH_URL указывает на порт 7100"
else
    echo "❌ ОШИБКА: AUTH_URL НЕ указывает на порт 7100!"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "================================"
if [ $ERRORS -eq 0 ]; then
    echo "✅ Все проверки пройдены!"
    echo "Авторизация настроена правильно (порт 7100)"
    exit 0
else
    echo "❌ Найдено ошибок: $ERRORS"
    echo "Исправь использование портов для авторизации!"
    echo ""
    echo "ПРАВИЛО: Авторизация ВСЕГДА через Auth Service (7100)"
    echo "  ✅ authClient.post('/api/auth/login', ...)"
    echo "  ✅ fetch('http://localhost:7100/api/auth/me', ...)"
    echo "  ❌ apiClient.post('/api/auth/login', ...)"
    echo "  ❌ fetch('http://localhost:8000/api/auth/me', ...)"
    exit 1
fi
