#!/bin/bash

echo "🔍 Отладка SSO"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Проверяем секреты
echo "1️⃣ Проверка JWT_SECRET:"
echo "───────────────────────────────────────────────────────────"
echo "Main Backend (.env):"
MAIN_SECRET=$(grep JWT_SECRET backend/.env | cut -d'=' -f2-)
echo "  $MAIN_SECRET"
echo "  Длина: ${#MAIN_SECRET}"
echo ""
echo "Admin Backend (.env):"
ADMIN_SECRET=$(grep JWT_SECRET admin/backend/.env | cut -d'=' -f2-)
echo "  $ADMIN_SECRET"
echo "  Длина: ${#ADMIN_SECRET}"
echo ""

if [ "$MAIN_SECRET" = "$ADMIN_SECRET" ]; then
  echo "✅ Секреты одинаковые (длина: ${#MAIN_SECRET})"
else
  echo "❌ ОШИБКА: Секреты разные!"
  echo "Main: '$MAIN_SECRET'"
  echo "Admin: '$ADMIN_SECRET'"
  exit 1
fi
echo ""

# Проверяем что backend запущен
echo "2️⃣ Проверка сервисов:"
echo "───────────────────────────────────────────────────────────"
echo "Main Backend (8000):"
MAIN_HEALTH=$(curl -s http://localhost:8000/api/health)
if [ -n "$MAIN_HEALTH" ]; then
  echo "  ✅ Запущен: $MAIN_HEALTH"
else
  echo "  ❌ Не запущен"
  exit 1
fi
echo ""

echo "Admin Backend (9000):"
ADMIN_HEALTH=$(curl -s http://localhost:9000/api/admin/health)
if [ -n "$ADMIN_HEALTH" ]; then
  echo "  ✅ Запущен: $ADMIN_HEALTH"
else
  echo "  ❌ Не запущен"
  exit 1
fi
echo ""

# Логинимся
echo "3️⃣ Вход в систему:"
echo "───────────────────────────────────────────────────────────"
echo "Отправка запроса на http://localhost:8000/api/auth/login"
echo "Email: anton@dvinyaninov.ru"
RESPONSE=$(curl -s -c /tmp/debug-cookies.txt -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"anton@dvinyaninov.ru","password":"dxG0BBG0"}')

echo "Ответ:"
echo "$RESPONSE" | jq '.'

if echo "$RESPONSE" | grep -q '"success":true'; then
  echo "✅ Вход успешен"
else
  echo "❌ Ошибка входа"
  exit 1
fi
echo ""

# Показываем cookie
echo "4️⃣ Cookie файл:"
echo "───────────────────────────────────────────────────────────"
cat /tmp/debug-cookies.txt
echo ""

# Извлекаем токен
TOKEN=$(cat /tmp/debug-cookies.txt | grep auth_token | awk '{print $7}')
echo "Токен:"
echo "  $TOKEN"
echo ""

# Декодируем payload токена
echo "Payload токена (декодированный):"
PAYLOAD=$(echo "$TOKEN" | cut -d'.' -f2)
# Добавляем padding если нужно
PADDING=$((4 - ${#PAYLOAD} % 4))
if [ $PADDING -ne 4 ]; then
  PAYLOAD="${PAYLOAD}$(printf '=%.0s' $(seq 1 $PADDING))"
fi
echo "$PAYLOAD" | base64 -d 2>/dev/null | jq '.' || echo "  (не удалось декодировать)"
echo ""

# Проверяем токен через Main Backend
echo "5️⃣ Проверка токена через Main Backend:"
echo "───────────────────────────────────────────────────────────"
echo "GET http://localhost:8000/api/auth/verify"
VERIFY_RESPONSE=$(curl -s -b /tmp/debug-cookies.txt http://localhost:8000/api/auth/verify)
echo "Ответ:"
echo "$VERIFY_RESPONSE" | jq '.'

if echo "$VERIFY_RESPONSE" | grep -q '"valid":true'; then
  echo "✅ Main Backend валидирует токен"
  
  # Проверяем наличие ролей
  if echo "$VERIFY_RESPONSE" | grep -q '"superadmin"'; then
    echo "✅ Роль superadmin найдена"
  else
    echo "❌ Роль superadmin отсутствует"
  fi
else
  echo "❌ Main Backend не валидирует токен"
fi
echo ""

# Проверяем доступ к админке
echo "6️⃣ Проверка доступа к админке:"
echo "───────────────────────────────────────────────────────────"
echo "GET http://localhost:9000/api/admin/auth/me"
echo "Cookie: auth_token=$TOKEN"
echo ""
ADMIN_RESPONSE=$(curl -s -b /tmp/debug-cookies.txt http://localhost:9000/api/admin/auth/me)
echo "Ответ:"
echo "$ADMIN_RESPONSE" | jq '.'

if echo "$ADMIN_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Доступ к админке разрешен"
else
  echo "❌ Доступ к админке запрещен"
  
  # Показываем ошибку
  ERROR=$(echo "$ADMIN_RESPONSE" | jq -r '.error')
  echo "Ошибка: $ERROR"
fi
echo ""

# Проверяем логи Admin Backend
echo "7️⃣ Проверка логов Admin Backend:"
echo "───────────────────────────────────────────────────────────"
echo "Последние 10 строк логов:"
echo "(Запустите Admin Backend вручную чтобы увидеть логи)"
echo ""

echo "═══════════════════════════════════════════════════════════"
echo "Тест завершен"
echo ""

rm -f /tmp/debug-cookies.txt
