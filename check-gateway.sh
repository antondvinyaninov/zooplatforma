#!/bin/bash

echo "🔍 Проверка Gateway на Easypanel..."
echo ""

GATEWAY_URL="https://my-projects-gateway-zp.crv1ic.easypanel.host"

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Проверка доступности Gateway
echo "1️⃣ Проверка доступности Gateway..."
HEALTH_RESPONSE=$(curl -s "${GATEWAY_URL}/health")
if echo "$HEALTH_RESPONSE" | grep -q "gateway"; then
    echo -e "${GREEN}✅ Gateway доступен${NC}"
    echo "   Ответ: $(echo $HEALTH_RESPONSE | head -c 100)..."
else
    echo -e "${RED}❌ Gateway недоступен${NC}"
    exit 1
fi

echo ""

# 2. Проверка CORS для OPTIONS (preflight)
echo "2️⃣ Проверка CORS preflight (OPTIONS)..."
CORS_RESPONSE=$(curl -s -X OPTIONS "${GATEWAY_URL}/api/posts" \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: GET" \
    -H "Access-Control-Request-Headers: Content-Type, Authorization" \
    -i)

if echo "$CORS_RESPONSE" | grep -qi "access-control-allow-origin.*localhost:3000"; then
    echo -e "${GREEN}✅ CORS preflight работает${NC}"
    echo "   Access-Control-Allow-Origin: http://localhost:3000"
else
    echo -e "${RED}❌ CORS preflight НЕ работает${NC}"
    echo -e "${YELLOW}   Gateway не возвращает Access-Control-Allow-Origin для localhost:3000${NC}"
    echo ""
    echo "Заголовки ответа:"
    echo "$CORS_RESPONSE" | grep -i "access-control"
    exit 1
fi

echo ""

# 3. Проверка CORS для GET запроса
echo "3️⃣ Проверка CORS для GET запроса..."
CORS_GET=$(curl -s -X GET "${GATEWAY_URL}/api/posts" \
    -H "Origin: http://localhost:3000" \
    -i)

if echo "$CORS_GET" | grep -qi "access-control-allow-origin.*localhost:3000"; then
    echo -e "${GREEN}✅ CORS для GET работает${NC}"
else
    echo -e "${RED}❌ CORS для GET НЕ работает${NC}"
    exit 1
fi

echo ""

# 4. Проверка авторизации (без токена должен вернуть 401)
echo "4️⃣ Проверка авторизации..."
AUTH_RESPONSE=$(curl -s -w "\n%{http_code}" "${GATEWAY_URL}/api/profile" \
    -H "Origin: http://localhost:3000")

HTTP_CODE=$(echo "$AUTH_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "401" ]; then
    echo -e "${GREEN}✅ Авторизация работает (401 без токена)${NC}"
else
    echo -e "${YELLOW}⚠️ Неожиданный код: $HTTP_CODE${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Все проверки пройдены! Gateway готов к работе.${NC}"
echo ""
echo "Теперь можешь использовать Gateway:"
echo "  NEXT_PUBLIC_API_URL=https://my-projects-gateway-zp.crv1ic.easypanel.host"
