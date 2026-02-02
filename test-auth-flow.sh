#!/bin/bash

echo "🔐 Testing Auth Flow"
echo "===================="
echo ""

# 1. Login to Auth Service
echo "1️⃣ Login to Auth Service (port 7100)..."
LOGIN_RESPONSE=$(curl -s -c cookies.txt -X POST http://localhost:7100/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

echo "Response: $LOGIN_RESPONSE"
echo ""

# Check if login successful
if echo "$LOGIN_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Login successful"
  
  # Extract token from response
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  echo "Token: ${TOKEN:0:50}..."
  echo ""
else
  echo "❌ Login failed"
  exit 1
fi

# 2. Check /api/auth/me with cookie
echo "2️⃣ Check /api/auth/me with cookie..."
ME_RESPONSE=$(curl -s -b cookies.txt http://localhost:7100/api/auth/me)
echo "Response: $ME_RESPONSE"
echo ""

if echo "$ME_RESPONSE" | grep -q '"success":true'; then
  echo "✅ /api/auth/me with cookie works"
  
  # Check if token is in response
  if echo "$ME_RESPONSE" | grep -q '"token"'; then
    echo "✅ Token is in /api/auth/me response"
  else
    echo "❌ Token is NOT in /api/auth/me response"
  fi
  echo ""
else
  echo "❌ /api/auth/me with cookie failed"
fi

# 3. Check /api/auth/me with Bearer token
echo "3️⃣ Check /api/auth/me with Bearer token..."
ME_TOKEN_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:7100/api/auth/me)
echo "Response: $ME_TOKEN_RESPONSE"
echo ""

if echo "$ME_TOKEN_RESPONSE" | grep -q '"success":true'; then
  echo "✅ /api/auth/me with Bearer token works"
else
  echo "❌ /api/auth/me with Bearer token failed"
fi

# 4. Test Main Backend with Bearer token
echo "4️⃣ Test Main Backend /api/friends with Bearer token..."
FRIENDS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/friends)
echo "Response: $FRIENDS_RESPONSE"
echo ""

if echo "$FRIENDS_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Main Backend /api/friends works with Bearer token"
else
  echo "❌ Main Backend /api/friends failed"
fi

# 5. Test Main Backend /api/notifications/unread
echo "5️⃣ Test Main Backend /api/notifications/unread with Bearer token..."
NOTIF_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/notifications/unread)
echo "Response: $NOTIF_RESPONSE"
echo ""

if echo "$NOTIF_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Main Backend /api/notifications/unread works"
else
  echo "❌ Main Backend /api/notifications/unread failed"
fi

# 6. Test Main Backend /api/organizations/my
echo "6️⃣ Test Main Backend /api/organizations/my with Bearer token..."
ORG_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/organizations/my)
echo "Response: $ORG_RESPONSE"
echo ""

if echo "$ORG_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Main Backend /api/organizations/my works"
else
  echo "❌ Main Backend /api/organizations/my failed"
fi

# Cleanup
rm -f cookies.txt

echo ""
echo "===================="
echo "🎉 Auth Flow Test Complete"
