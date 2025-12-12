#!/bin/bash

# Тест SSO (Single Sign-On) системы
# Проверяет, что авторизация работает между микросервисами

echo "🔐 Тестирование SSO системы..."
echo ""

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# URLs
MAIN_API="http://localhost:8080"
ADMIN_API="http://localhost:8081"
COOKIE_FILE="/tmp/sso-test-cookies.txt"

# Очищаем старые cookies
rm -f "$COOKIE_FILE"

echo "📋 Тест 1: Вход на основном сайте"
echo "=================================="

# Регистрируем тестового пользователя (или используем существующего)
LOGIN_RESPONSE=$(curl -s -c "$COOKIE_FILE" -X POST "$MAIN_API/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

if echo "$LOGIN_RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✅ Вход успешен${NC}"
  echo "$LOGIN_RESPONSE" | jq '.'
else
  echo -e "${RED}❌ Ошибка входа${NC}"
  echo "$LOGIN_RESPONSE" | jq '.'
  
  echo ""
  echo "Попытка регистрации нового пользователя..."
  REGISTER_RESPONSE=$(curl -s -c "$COOKIE_FILE" -X POST "$MAIN_API/api/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test User","email":"test@example.com","password":"password123"}')
  
  if echo "$REGISTER_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Регистрация успешна${NC}"
  else
    echo -e "${RED}❌ Ошибка регистрации${NC}"
    exit 1
  fi
fi

echo ""
echo "📋 Тест 2: Проверка токена на основном сайте"
echo "============================================="

ME_RESPONSE=$(curl -s -b "$COOKIE_FILE" "$MAIN_API/api/auth/me")

if echo "$ME_RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✅ Токен валиден${NC}"
  echo "$ME_RESPONSE" | jq '.'
else
  echo -e "${RED}❌ Токен невалиден${NC}"
  exit 1
fi

echo ""
echo "📋 Тест 3: Проверка токена через /api/auth/verify"
echo "=================================================="

VERIFY_RESPONSE=$(curl -s -b "$COOKIE_FILE" "$MAIN_API/api/auth/verify")

if echo "$VERIFY_RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✅ Endpoint /api/auth/verify работает${NC}"
  echo "$VERIFY_RESPONSE" | jq '.'
  
  # Проверяем наличие roles
  if echo "$VERIFY_RESPONSE" | grep -q '"roles"'; then
    echo -e "${GREEN}✅ Роли присутствуют в токене${NC}"
  else
    echo -e "${YELLOW}⚠️  Роли отсутствуют в токене (старый токен?)${NC}"
  fi
else
  echo -e "${RED}❌ Endpoint /api/auth/verify не работает${NC}"
  exit 1
fi

echo ""
echo "📋 Тест 4: Доступ к админ-панели (без прав)"
echo "============================================"

ADMIN_ME_RESPONSE=$(curl -s -b "$COOKIE_FILE" "$ADMIN_API/api/admin/auth/me")

if echo "$ADMIN_ME_RESPONSE" | grep -q '"success":false'; then
  echo -e "${GREEN}✅ Доступ запрещен (ожидаемо, нет прав админа)${NC}"
  echo "$ADMIN_ME_RESPONSE" | jq '.'
else
  echo -e "${YELLOW}⚠️  Неожиданный ответ${NC}"
  echo "$ADMIN_ME_RESPONSE" | jq '.'
fi

echo ""
echo "📋 Тест 5: Проверка наличия суперадмина"
echo "========================================"

# Получаем ID пользователя из токена
USER_ID=$(echo "$ME_RESPONSE" | jq -r '.data.id')

if [ "$USER_ID" != "null" ] && [ -n "$USER_ID" ]; then
  echo "User ID: $USER_ID"
  
  # Проверяем, является ли пользователь суперадмином
  ADMIN_CHECK=$(sqlite3 database/data.db "SELECT COUNT(*) FROM admins WHERE user_id = $USER_ID AND role = 'superadmin';")
  
  if [ "$ADMIN_CHECK" -gt 0 ]; then
    echo -e "${GREEN}✅ Пользователь является суперадмином${NC}"
    
    echo ""
    echo "📋 Тест 6: Доступ к админ-панели (с правами)"
    echo "============================================"
    
    # Нужно перелогиниться, чтобы получить новый токен с ролями
    rm -f "$COOKIE_FILE"
    curl -s -c "$COOKIE_FILE" -X POST "$MAIN_API/api/auth/login" \
      -H "Content-Type: application/json" \
      -d '{"email":"test@example.com","password":"password123"}' > /dev/null
    
    ADMIN_ME_RESPONSE=$(curl -s -b "$COOKIE_FILE" "$ADMIN_API/api/admin/auth/me")
    
    if echo "$ADMIN_ME_RESPONSE" | grep -q '"success":true'; then
      echo -e "${GREEN}✅ Доступ к админ-панели разрешен${NC}"
      echo "$ADMIN_ME_RESPONSE" | jq '.'
    else
      echo -e "${RED}❌ Доступ к админ-панели запрещен${NC}"
      echo "$ADMIN_ME_RESPONSE" | jq '.'
    fi
  else
    echo -e "${YELLOW}⚠️  Пользователь не является суперадмином${NC}"
    echo ""
    echo "Для назначения суперадмином выполните:"
    echo "  cd admin/backend"
    echo "  ./create-superadmin.sh $USER_ID"
  fi
else
  echo -e "${RED}❌ Не удалось получить ID пользователя${NC}"
fi

echo ""
echo "📋 Тест 7: Выход из системы"
echo "==========================="

LOGOUT_RESPONSE=$(curl -s -b "$COOKIE_FILE" -c "$COOKIE_FILE" -X POST "$MAIN_API/api/auth/logout")

if echo "$LOGOUT_RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✅ Выход успешен${NC}"
else
  echo -e "${RED}❌ Ошибка выхода${NC}"
fi

# Проверяем, что токен больше не работает
ME_AFTER_LOGOUT=$(curl -s -b "$COOKIE_FILE" "$MAIN_API/api/auth/me")

if echo "$ME_AFTER_LOGOUT" | grep -q '"success":false'; then
  echo -e "${GREEN}✅ Токен удален${NC}"
else
  echo -e "${RED}❌ Токен все еще работает${NC}"
fi

# Очищаем
rm -f "$COOKIE_FILE"

echo ""
echo "=================================="
echo "🎉 Тестирование SSO завершено!"
echo "=================================="
