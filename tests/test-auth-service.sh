#!/bin/bash

# Тест Auth Service

echo "🧪 Testing Auth Service..."
echo ""

AUTH_URL="http://localhost:7100"

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Health Check
echo "1️⃣ Testing Health Check..."
response=$(curl -s "$AUTH_URL/api/health")
if echo "$response" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Health check passed${NC}"
    echo "Response: $response"
else
    echo -e "${RED}❌ Health check failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 2. Регистрация нового пользователя
echo "2️⃣ Testing Registration..."
RANDOM_EMAIL="test_$(date +%s)@example.com"
REGISTER_DATA='{
  "email": "'$RANDOM_EMAIL'",
  "password": "password123",
  "name": "Test",
  "last_name": "User"
}'

response=$(curl -s -X POST "$AUTH_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "$REGISTER_DATA")

if echo "$response" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Registration passed${NC}"
    TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    USER_ID=$(echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    echo "Email: $RANDOM_EMAIL"
    echo "User ID: $USER_ID"
    echo "Token: ${TOKEN:0:50}..."
else
    echo -e "${RED}❌ Registration failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 3. Проверка токена
echo "3️⃣ Testing Token Verification..."
response=$(curl -s -X GET "$AUTH_URL/api/auth/verify" \
  -H "Authorization: Bearer $TOKEN")

if echo "$response" | grep -q '"valid":true'; then
    echo -e "${GREEN}✅ Token verification passed${NC}"
    echo "Response: $response"
else
    echo -e "${RED}❌ Token verification failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 4. Получение информации о пользователе
echo "4️⃣ Testing Get Me..."
response=$(curl -s -X GET "$AUTH_URL/api/auth/me" \
  -H "Authorization: Bearer $TOKEN")

if echo "$response" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Get Me passed${NC}"
    echo "Response: $response"
else
    echo -e "${RED}❌ Get Me failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 5. Логин с теми же credentials
echo "5️⃣ Testing Login..."
LOGIN_DATA='{
  "email": "'$RANDOM_EMAIL'",
  "password": "password123"
}'

response=$(curl -s -X POST "$AUTH_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "$LOGIN_DATA")

if echo "$response" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Login passed${NC}"
    NEW_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "New Token: ${NEW_TOKEN:0:50}..."
else
    echo -e "${RED}❌ Login failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 6. Проверка неправильного пароля
echo "6️⃣ Testing Wrong Password..."
WRONG_LOGIN_DATA='{
  "email": "'$RANDOM_EMAIL'",
  "password": "wrongpassword"
}'

response=$(curl -s -X POST "$AUTH_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "$WRONG_LOGIN_DATA")

if echo "$response" | grep -q '"success":false'; then
    echo -e "${GREEN}✅ Wrong password correctly rejected${NC}"
else
    echo -e "${RED}❌ Wrong password test failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 7. Проверка несуществующего пользователя
echo "7️⃣ Testing Non-existent User..."
NONEXIST_LOGIN_DATA='{
  "email": "nonexistent@example.com",
  "password": "password123"
}'

response=$(curl -s -X POST "$AUTH_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "$NONEXIST_LOGIN_DATA")

if echo "$response" | grep -q '"success":false'; then
    echo -e "${GREEN}✅ Non-existent user correctly rejected${NC}"
else
    echo -e "${RED}❌ Non-existent user test failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# 8. Проверка невалидного токена
echo "8️⃣ Testing Invalid Token..."
response=$(curl -s -X GET "$AUTH_URL/api/auth/verify" \
  -H "Authorization: Bearer invalid_token_here")

if echo "$response" | grep -q '"valid":false'; then
    echo -e "${GREEN}✅ Invalid token correctly rejected${NC}"
else
    echo -e "${RED}❌ Invalid token test failed${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# Итоги
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 All tests passed!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Test User:"
echo "  Email: $RANDOM_EMAIL"
echo "  Password: password123"
echo "  User ID: $USER_ID"
echo "  Token: ${TOKEN:0:50}..."
echo ""
echo "You can use this token to test other services!"
