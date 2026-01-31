#!/bin/bash

echo "🧪 Testing Centralized Auth System..."
echo ""

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

PASSED=0
FAILED=0

check_test() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED${NC}: $2"
        ((PASSED++))
    else
        echo -e "${RED}❌ FAILED${NC}: $2"
        ((FAILED++))
    fi
}

# 1. Auth Service health
echo "1️⃣ Testing Auth Service..."
RESPONSE=$(curl -s http://localhost:7100/api/health)
if echo "$RESPONSE" | grep -q "ok"; then
    check_test 0 "Auth Service is running"
else
    check_test 1 "Auth Service is NOT running"
    exit 1
fi
echo ""

# 2. Register user
echo "2️⃣ Testing registration..."
TIMESTAMP=$(date +%s)
TEST_EMAIL="test_${TIMESTAMP}@example.com"

REGISTER_RESPONSE=$(curl -s -X POST http://localhost:7100/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"password123\",\"name\":\"Test\"}")

if echo "$REGISTER_RESPONSE" | grep -q "success.*true"; then
    check_test 0 "User registration"
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
else
    check_test 1 "User registration"
    exit 1
fi
echo ""

# 3. Verify token
echo "3️⃣ Testing token verification..."
VERIFY=$(curl -s http://localhost:7100/api/auth/verify -H "Authorization: Bearer $TOKEN")
if echo "$VERIFY" | grep -q "valid.*true"; then
    check_test 0 "Token verification"
else
    check_test 1 "Token verification"
fi
echo ""

# 4. Test protected endpoint
echo "4️⃣ Testing protected endpoint..."
PROTECTED=$(curl -s -w "%{http_code}" http://localhost:8100/api/pets -H "Authorization: Bearer $TOKEN")
HTTP_CODE="${PROTECTED: -3}"
if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
    check_test 0 "Protected endpoint with token"
else
    check_test 1 "Protected endpoint (got $HTTP_CODE)"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Passed: $PASSED${NC}"
echo -e "${RED}❌ Failed: $FAILED${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    exit 1
fi
