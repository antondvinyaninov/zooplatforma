#!/bin/bash

echo "🔐 Проверка доступа к админ-панели"
echo ""
echo "Инструкция:"
echo "1. Откройте DevTools в браузере (F12)"
echo "2. Перейдите на вкладку Application → Cookies → http://localhost:3000"
echo "3. Найдите cookie 'auth_token' и скопируйте его значение"
echo "4. Вставьте значение ниже и нажмите Enter"
echo ""
read -p "Введите значение auth_token: " TOKEN

if [ -z "$TOKEN" ]; then
  echo "❌ Токен не введен"
  exit 1
fi

echo ""
echo "Проверяю доступ..."
echo ""

# Проверка /api/admin/auth/me
echo "📋 GET /api/admin/auth/me"
RESPONSE=$(curl -s -H "Cookie: auth_token=$TOKEN" http://localhost:8083/api/admin/auth/me)
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

echo ""
echo "📋 GET /api/admin/users"
RESPONSE=$(curl -s -H "Cookie: auth_token=$TOKEN" http://localhost:8083/api/admin/users)
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

echo ""
echo "📋 GET /api/admin/stats/overview"
RESPONSE=$(curl -s -H "Cookie: auth_token=$TOKEN" http://localhost:8083/api/admin/stats/overview)
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
