# Troubleshooting Guide: Auth & CORS

## Overview

This guide helps diagnose and fix common authentication and CORS issues in the Main Service.

---

## Quick Diagnostic Commands

### Check if services are running

```bash
# Gateway
curl http://localhost/health

# Backend
curl http://localhost:8000/api/health

# Frontend
curl http://localhost:3000
```

### Check JWT token

```bash
# Get token from localStorage (in browser console)
localStorage.getItem('auth_token')

# Decode JWT (use jwt.io or)
echo "YOUR_TOKEN" | cut -d'.' -f2 | base64 -d | jq
```

### Check CORS headers

```bash
# Preflight request
curl -X OPTIONS http://localhost/api/posts \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  -v

# Should see:
# Access-Control-Allow-Origin: http://localhost:3000
# Access-Control-Allow-Credentials: true
```

### Check Gateway headers

```bash
# Make authenticated request
curl http://localhost/api/posts \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Origin: http://localhost:3000" \
  -v

# Should see X-User-* headers in backend logs
```

---

## Common Errors

### Error 1: 401 Unauthorized

#### Symptoms

```
GET http://localhost:8000/api/posts 401 (Unauthorized)
```

Browser console shows:
```
❌ Failed to load posts: Unauthorized
```

#### Possible Causes

**Cause 1: No token sent**

```typescript
// ❌ НЕПРАВИЛЬНО: забыли credentials
fetch('/api/posts')

// ✅ ПРАВИЛЬНО: используй apiClient
apiClient.get('/api/posts')
```

**Cause 2: Token expired**

JWT tokens expire after 7 days. Check token expiration:

```javascript
// In browser console
const token = localStorage.getItem('auth_token');
const payload = JSON.parse(atob(token.split('.')[1]));
console.log('Expires:', new Date(payload.exp * 1000));
console.log('Now:', new Date());
```

**Solution:** Login again to get new token

**Cause 3: JWT_SECRET mismatch**

```bash
# Check Gateway
cat gateway/.env | grep JWT_SECRET

# Check Backend
cat backend/.env | grep JWT_SECRET

# Must be the same!
```

**Solution:** Update JWT_SECRET to match

**Cause 4: Cookie not sent**

```typescript
// ❌ НЕПРАВИЛЬНО: забыли credentials
fetch('/api/posts', {
  headers: { Authorization: 'Bearer TOKEN' }
})

// ✅ ПРАВИЛЬНО: добавь credentials
fetch('/api/posts', {
  headers: { Authorization: 'Bearer TOKEN' },
  credentials: 'include'  // ← Important!
})
```

#### Diagnostic Steps

1. **Check if token exists:**
   ```javascript
   // Browser console
   localStorage.getItem('auth_token')
   // Should return token string
   ```

2. **Check if token is sent:**
   ```
   DevTools → Network → Select request → Request Headers
   Should see: Authorization: Bearer ...
   Or: Cookie: auth_token=...
   ```

3. **Check backend logs:**
   ```bash
   tail -f /tmp/main-backend.log
   # Should see: ✅ User authenticated: id=123
   # If see: ⚠️ No token found - token not sent
   # If see: ❌ Invalid token - token invalid
   ```

4. **Check JWT_SECRET:**
   ```bash
   # Gateway and Backend must have same secret
   diff <(cat gateway/.env | grep JWT_SECRET) \
        <(cat backend/.env | grep JWT_SECRET)
   # Should be no difference
   ```

#### Solutions

**Solution 1: Use apiClient**

```typescript
// ❌ Before
const response = await fetch('/api/posts');

// ✅ After
import { apiClient } from '@/lib/api';
const response = await apiClient.get('/api/posts');
```

**Solution 2: Login again**

```typescript
// Token expired - redirect to login
if (response.error === 'Unauthorized') {
  router.push('/auth');
}
```

**Solution 3: Fix JWT_SECRET**

```bash
# Copy secret from Gateway to Backend
JWT_SECRET=$(cat gateway/.env | grep JWT_SECRET | cut -d'=' -f2)
echo "JWT_SECRET=$JWT_SECRET" >> backend/.env
```

---

### Error 2: CORS Policy Error

#### Symptoms

```
Access to fetch at 'http://localhost:8000/api/posts' from origin 
'http://localhost:3000' has been blocked by CORS policy: 
No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

Or:

```
Access to fetch at 'http://localhost:8000/api/posts' from origin 
'http://localhost:3000' has been blocked by CORS policy: 
The 'Access-Control-Allow-Origin' header contains multiple values 
'http://localhost:3000, http://localhost:3000', but only one is allowed.
```

#### Possible Causes

**Cause 1: Backend also sets CORS headers**

```go
// ❌ НЕПРАВИЛЬНО: Backend устанавливает CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        next(w, r)
    }
}
```

**Result:** Gateway и Backend оба устанавливают CORS → дублирование

**Solution:** Убери CORS из Backend, оставь только в Gateway

**Cause 2: Origin not in allowedOrigins**

```go
// Gateway middleware
allowedOrigins := map[string]bool{
    "http://localhost:3000": true,
    // "http://localhost:3001": true,  // ← Забыли добавить!
}
```

**Solution:** Добавь origin в allowedOrigins

**Cause 3: Preflight request blocked**

```go
// ❌ НЕПРАВИЛЬНО: OPTIONS не обрабатывается
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusForbidden)  // ← Блокирует preflight!
    return
}
```

**Solution:** Верни 200 OK для OPTIONS

#### Diagnostic Steps

1. **Check if CORS headers are duplicated:**
   ```bash
   curl -H "Origin: http://localhost:3000" \
        http://localhost:8000/api/posts -v 2>&1 | \
        grep "Access-Control-Allow-Origin"
   
   # Should see ONLY ONE line:
   # Access-Control-Allow-Origin: http://localhost:3000
   
   # If see TWO lines - headers duplicated!
   ```

2. **Check Gateway allowedOrigins:**
   ```bash
   cat gateway/middleware/cors.go | grep -A 5 "allowedOrigins"
   # Should see your origin in the list
   ```

3. **Check preflight response:**
   ```bash
   curl -X OPTIONS http://localhost/api/posts \
     -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -v
   
   # Should return 200 OK
   # Should have Access-Control-* headers
   ```

4. **Check Gateway logs:**
   ```bash
   # Gateway should log CORS decisions
   tail -f gateway/logs/gateway.log
   
   # Should see:
   # ✅ CORS: Allowed origin http://localhost:3000
   # Or:
   # ⚠️ CORS: Blocked origin http://unknown:3000
   ```

#### Solutions

**Solution 1: Remove CORS from Backend**

```go
// backend/main.go

// ❌ Remove this:
http.HandleFunc("/api/posts", enableCORS(handlers.HandlePosts))

// ✅ Use this:
http.HandleFunc("/api/posts", handlers.HandlePosts)
```

**Solution 2: Add origin to Gateway**

```go
// gateway/middleware/cors.go

allowedOrigins := map[string]bool{
    "http://localhost:3000": true,
    "http://localhost:3001": true,  // ← Add new origin
    "https://your-frontend.com": true,
}
```

**Solution 3: Fix preflight handling**

```go
// gateway/middleware/cors.go

if allowedOrigins[origin] {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
    
    // ✅ Return 200 OK for preflight
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
}
```

---

### Error 3: Failed to fetch

#### Symptoms

```
TypeError: Failed to fetch
    at apiClient.get (api.ts:84)
```

No network request appears in DevTools Network tab.

#### Possible Causes

**Cause 1: Service not running**

```bash
# Check if Backend is running
curl http://localhost:8000/api/health
# If fails - Backend not running
```

**Solution:** Start Backend

```bash
cd backend
go run main.go
```

**Cause 2: Wrong URL**

```typescript
// ❌ НЕПРАВИЛЬНО: хардкод URL
const response = await fetch('http://localhost:9999/api/posts');
// Port 9999 doesn't exist!

// ✅ ПРАВИЛЬНО: используй apiClient
const response = await apiClient.get('/api/posts');
```

**Cause 3: Network error**

Browser can't reach the server (firewall, VPN, etc.)

**Solution:** Check network connectivity

```bash
# Ping localhost
ping localhost

# Check if port is open
lsof -i :8000
```

**Cause 4: HTTPS mixed content**

```
Mixed Content: The page at 'https://frontend.com' was loaded over HTTPS, 
but requested an insecure resource 'http://backend.com/api/posts'. 
This request has been blocked.
```

**Solution:** Use HTTPS for all requests in production

```env
# .env.production
NEXT_PUBLIC_API_URL=https://gateway.com  # ← HTTPS!
```

#### Diagnostic Steps

1. **Check if service is running:**
   ```bash
   # Backend
   curl http://localhost:8000/api/health
   
   # Gateway
   curl http://localhost/health
   ```

2. **Check browser console:**
   ```
   DevTools → Console
   Look for error details
   ```

3. **Check Network tab:**
   ```
   DevTools → Network
   - If request appears → server responded (check status code)
   - If request doesn't appear → request blocked before sending
   ```

4. **Check CORS:**
   ```bash
   # If request blocked by CORS, you'll see CORS error
   # See "Error 2: CORS Policy Error" above
   ```

#### Solutions

**Solution 1: Start services**

```bash
# Start all services
./run

# Or individually:
cd backend && go run main.go &
cd frontend && npm run dev &
```

**Solution 2: Use apiClient**

```typescript
// ❌ Before
const response = await fetch('http://localhost:8000/api/posts');

// ✅ After
import { apiClient } from '@/lib/api';
const response = await apiClient.get('/api/posts');
```

**Solution 3: Check .env**

```bash
# Development
cat frontend/.env.local
# Should have: NEXT_PUBLIC_API_URL=http://localhost:8000

# Production
cat frontend/.env.production
# Should have: NEXT_PUBLIC_API_URL= (empty or HTTPS URL)
```

---

### Error 4: Cookie not set after login

#### Symptoms

- Login succeeds (200 OK)
- But subsequent requests return 401 Unauthorized
- Cookie not visible in DevTools → Application → Cookies

#### Possible Causes

**Cause 1: Cookie Domain mismatch**

```go
// ❌ НЕПРАВИЛЬНО: Domain не подходит
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Domain: "example.com",  // ← Не подходит для localhost!
})
```

**Solution:** Use correct Domain

```go
// ✅ Development
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Domain: "localhost",  // ← Работает для всех портов localhost
})

// ✅ Production
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Domain: ".easypanel.host",  // ← Работает для всех поддоменов
})
```

**Cause 2: SameSite=Strict**

```go
// ❌ НЕПРАВИЛЬНО: Strict блокирует cross-site
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    SameSite: http.SameSiteStrictMode,  // ← Слишком строго!
})
```

**Solution:** Use SameSite=Lax

```go
// ✅ ПРАВИЛЬНО
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    SameSite: http.SameSiteLaxMode,  // ← Разрешает GET запросы
})
```

**Cause 3: Secure=true on HTTP**

```go
// ❌ НЕПРАВИЛЬНО: Secure требует HTTPS
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Secure: true,  // ← Не работает на http://localhost
})
```

**Solution:** Secure=false для development

```go
// ✅ Development
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Secure: false,  // ← HTTP OK
})

// ✅ Production
http.SetCookie(w, &http.Cookie{
    Name:   "auth_token",
    Value:  token,
    Secure: true,  // ← HTTPS required
})
```

#### Diagnostic Steps

1. **Check if cookie is set:**
   ```
   DevTools → Application → Cookies → localhost
   Should see: auth_token
   ```

2. **Check cookie attributes:**
   ```
   Domain: localhost (or .easypanel.host)
   Path: /
   HttpOnly: true
   Secure: false (dev) or true (prod)
   SameSite: Lax
   ```

3. **Check Set-Cookie header:**
   ```
   DevTools → Network → Select login request → Response Headers
   Should see: Set-Cookie: auth_token=...; Domain=localhost; ...
   ```

4. **Check backend logs:**
   ```bash
   tail -f /tmp/main-backend.log
   # Should see: ✅ Cookie set: auth_token
   ```

#### Solutions

**Solution 1: Fix cookie Domain**

```go
// backend/handlers/auth.go

http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    Path:     "/",
    Domain:   "localhost",  // ✅ For development
    HttpOnly: true,
    Secure:   false,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400 * 7,
})
```

**Solution 2: Check credentials**

```typescript
// frontend/lib/api.ts

async get<T>(endpoint: string): Promise<ApiResponse<T>> {
  const response = await fetch(url, {
    method: 'GET',
    headers,
    credentials: 'include',  // ✅ Important!
  });
}
```

---

### Error 5: X-User-ID header not found

#### Symptoms

Backend logs show:
```
⚠️ No X-User-ID header from Gateway
```

All requests return 401 Unauthorized in production.

#### Possible Causes

**Cause 1: Not using Gateway**

Frontend makes requests directly to Backend, bypassing Gateway.

```typescript
// ❌ НЕПРАВИЛЬНО: прямой запрос к Backend
const response = await fetch('http://backend:8000/api/posts');
// Gateway не добавил X-User-* заголовки!

// ✅ ПРАВИЛЬНО: через Gateway
const response = await fetch('/api/posts');
// Gateway добавит X-User-* заголовки
```

**Cause 2: Gateway not adding headers**

Gateway middleware not configured correctly.

**Cause 3: Using DevAuthMiddleware in production**

```go
// ❌ НЕПРАВИЛЬНО: DevAuthMiddleware в production
http.HandleFunc("/api/posts", middleware.DevAuthMiddleware(handlers.HandlePosts))

// ✅ ПРАВИЛЬНО: AuthMiddleware в production
http.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.HandlePosts))
```

#### Diagnostic Steps

1. **Check if request goes through Gateway:**
   ```bash
   # Check Gateway logs
   tail -f gateway/logs/gateway.log
   # Should see: → Proxying to Main Backend: GET /api/posts
   ```

2. **Check if Gateway adds headers:**
   ```bash
   # Make request through Gateway
   curl http://localhost/api/posts \
     -H "Authorization: Bearer TOKEN" \
     -v
   
   # Check backend logs - should see X-User-ID
   tail -f /tmp/main-backend.log
   # Should see: ✅ User from Gateway: id=123
   ```

3. **Check middleware:**
   ```bash
   # Check which middleware is used
   cat backend/main.go | grep "HandlePosts"
   # Should see: middleware.AuthMiddleware (not DevAuthMiddleware)
   ```

#### Solutions

**Solution 1: Use relative paths**

```typescript
// frontend/lib/api.ts

// ✅ Production: пустая строка = относительные пути
const API_URL = process.env.NEXT_PUBLIC_API_URL || '';

// Requests will go to /api/posts (через Gateway)
```

**Solution 2: Fix Gateway middleware**

```go
// gateway/middleware/auth.go

func AuthMiddleware(next http.Handler) http.Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... validate JWT ...
        
        // ✅ Add headers
        r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
        r.Header.Set("X-User-Email", email)
        r.Header.Set("X-User-Role", role)
        
        next.ServeHTTP(w, r)
    }
}
```

**Solution 3: Use correct middleware**

```go
// backend/main.go

// ✅ Production: AuthMiddleware (reads X-User-* headers)
http.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.HandlePosts))

// ✅ Development: DevAuthMiddleware (validates JWT locally)
// http.HandleFunc("/api/posts", middleware.DevAuthMiddleware(handlers.HandlePosts))
```

---

## Verification Checklist

### After fixing authentication issues:

- [ ] Login works (200 OK)
- [ ] Cookie is set (visible in DevTools)
- [ ] Token is valid (not expired)
- [ ] JWT_SECRET matches everywhere
- [ ] Subsequent requests work (not 401)
- [ ] User data appears in backend logs

### After fixing CORS issues:

- [ ] Preflight (OPTIONS) returns 200 OK
- [ ] Access-Control-Allow-Origin header present
- [ ] Access-Control-Allow-Origin NOT duplicated
- [ ] Origin in allowedOrigins list
- [ ] Credentials work (cookies sent)
- [ ] No CORS errors in browser console

### After fixing network issues:

- [ ] Services are running
- [ ] Ports are correct
- [ ] URLs are correct (relative in prod)
- [ ] apiClient is used (not direct fetch)
- [ ] .env files configured correctly
- [ ] Requests appear in Network tab

---

## Getting Help

### Collect diagnostic information:

```bash
# 1. Check service status
curl http://localhost/health
curl http://localhost:8000/api/health

# 2. Check environment
cat frontend/.env.local
cat backend/.env | grep -E "(JWT_SECRET|DATABASE_URL)"

# 3. Check logs
tail -50 /tmp/main-backend.log
tail -50 gateway/logs/gateway.log

# 4. Check CORS
curl -X OPTIONS http://localhost/api/posts \
  -H "Origin: http://localhost:3000" \
  -v 2>&1 | grep "Access-Control"

# 5. Test authentication
curl http://localhost/api/auth/login \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}' \
  -c cookies.txt

curl http://localhost/api/posts \
  -b cookies.txt
```

### Include in bug report:

1. Error message (exact text)
2. Browser console output
3. Network tab screenshot
4. Backend logs
5. Gateway logs (if using Gateway)
6. Environment (development/production)
7. Steps to reproduce

---

## Prevention

### Best Practices:

1. **Always use apiClient** - не используй прямой fetch
2. **Check JWT_SECRET** - должен совпадать везде
3. **Use relative paths in production** - NEXT_PUBLIC_API_URL=
4. **Gateway manages CORS** - Backend не устанавливает CORS
5. **Test both environments** - development и production
6. **Monitor logs** - следи за ошибками в логах
7. **Use correct middleware** - DevAuthMiddleware (dev), AuthMiddleware (prod)

### Common Mistakes to Avoid:

- ❌ Хардкод URL в fetch
- ❌ Забыть credentials: 'include'
- ❌ Разные JWT_SECRET
- ❌ Backend устанавливает CORS
- ❌ Использовать DevAuthMiddleware в production
- ❌ Забыть добавить origin в allowedOrigins
- ❌ Secure=true на HTTP

---

**Last Updated:** 03.02.2026
