# Auth & CORS Architecture Specification

> **Status:** ✅ Completed  
> **Version:** 1.0.0  
> **Last Updated:** 03.02.2026

## Overview

This specification documents the authentication and CORS architecture for the Main Service in ZooPlatforma. The system uses a microservices architecture with an API Gateway as the single entry point, handling authentication and CORS for all services.

## Quick Links

- 📋 [Requirements](./requirements.md) - User stories and acceptance criteria
- 🏗️ [Design](./design.md) - Architecture diagrams and detailed design
- 🔧 [Troubleshooting](./troubleshooting.md) - Common issues and solutions
- ✅ [Tasks](./tasks.md) - Implementation tasks and checklist

## Key Concepts

### Authentication Flow

**Development:**
```
Frontend → Backend (DevAuthMiddleware validates JWT)
```

**Production:**
```
Frontend → Gateway (validates JWT, adds X-User-* headers) → Backend (reads headers)
```

### CORS Management

**Gateway is the single source of CORS headers:**
- ✅ Gateway sets Access-Control-* headers
- ✅ Backend does NOT set CORS headers
- ✅ Gateway filters CORS headers from backend responses

### API Client

**Always use apiClient for API requests:**
```typescript
import { apiClient } from '@/lib/api';

// ✅ Correct
const response = await apiClient.get('/api/posts');

// ❌ Wrong
const response = await fetch('http://localhost:8000/api/posts');
```

## Architecture Diagrams

### Development Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Development Mode                          │
└─────────────────────────────────────────────────────────────┘

Frontend (localhost:3000)
    ↓ fetch('http://localhost:8000/api/posts')
Backend (localhost:8000)
    ↓ DevAuthMiddleware validates JWT
Handler
    ↓ userID from context
PostgreSQL
```

### Production Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Mode                           │
└─────────────────────────────────────────────────────────────┘

Frontend (Easypanel)
    ↓ fetch('/api/posts')  // Relative path
Gateway (Easypanel)
    ↓ Validates JWT, adds X-User-* headers
Backend (Easypanel)
    ↓ Reads X-User-* headers
Handler
    ↓ userID from context
PostgreSQL
```

## Quick Start

### For Developers

1. **Read the Design Document:**
   ```bash
   cat .kiro/specs/auth-cors-architecture/design.md
   ```

2. **Understand the authentication flow:**
   - Development: Backend validates JWT locally
   - Production: Gateway validates JWT, Backend trusts headers

3. **Always use apiClient:**
   ```typescript
   import { apiClient } from '@/lib/api';
   const response = await apiClient.get('/api/posts');
   ```

4. **Never set CORS in Backend:**
   - Gateway manages all CORS
   - Backend should not have Access-Control-* headers

### For Troubleshooting

1. **Check the Troubleshooting Guide:**
   ```bash
   cat .kiro/specs/auth-cors-architecture/troubleshooting.md
   ```

2. **Common issues:**
   - 401 Unauthorized → Check JWT_SECRET, token expiration
   - CORS errors → Check Gateway allowedOrigins, remove Backend CORS
   - Failed to fetch → Check services running, check URLs

3. **Diagnostic commands:**
   ```bash
   # Check services
   curl http://localhost/health
   curl http://localhost:8000/api/health
   
   # Check CORS
   curl -X OPTIONS http://localhost/api/posts \
     -H "Origin: http://localhost:3000" -v
   
   # Check logs
   tail -f /tmp/main-backend.log
   ```

## Document Structure

### 1. Requirements Document

**Purpose:** Defines what needs to be documented

**Contents:**
- User stories for developers
- Acceptance criteria
- Requirements for documentation, diagrams, troubleshooting

**When to read:** When planning documentation updates

### 2. Design Document

**Purpose:** Describes how the system works

**Contents:**
- Architecture diagrams (development vs production)
- Authentication flow (login, token validation)
- CORS architecture (why Gateway manages CORS)
- API Client implementation
- Middleware architecture
- Security considerations

**When to read:** 
- When learning the system
- When implementing new features
- When debugging authentication/CORS issues

### 3. Troubleshooting Guide

**Purpose:** Helps fix common problems

**Contents:**
- Common errors (401, CORS, Failed to fetch)
- Diagnostic steps for each error
- Solutions with code examples
- Verification checklist
- Prevention best practices

**When to read:**
- When encountering errors
- When debugging issues
- Before deploying to production

### 4. Tasks Document

**Purpose:** Tracks implementation progress

**Contents:**
- Completed tasks (API Client, Middleware, CORS)
- Testing checklist
- Maintenance tasks
- Future improvements
- Rollback plan

**When to read:**
- When reviewing what was implemented
- When planning maintenance
- When considering improvements

## Key Principles

### 1. Gateway is the Single Entry Point

**All requests go through Gateway in production:**
- ✅ Gateway validates JWT tokens
- ✅ Gateway adds X-User-* headers
- ✅ Gateway manages CORS
- ✅ Backend trusts Gateway headers

### 2. Use apiClient for All API Requests

**Never use direct fetch:**
```typescript
// ❌ Wrong
const response = await fetch('http://localhost:8000/api/posts');

// ✅ Correct
import { apiClient } from '@/lib/api';
const response = await apiClient.get('/api/posts');
```

**Why:**
- Automatic URL configuration (dev/prod)
- Automatic credentials: 'include'
- Automatic Authorization header
- Unified error handling

### 3. Gateway Manages CORS

**Backend should NOT set CORS headers:**
```go
// ❌ Wrong - Backend sets CORS
w.Header().Set("Access-Control-Allow-Origin", "*")

// ✅ Correct - Backend doesn't set CORS
// Gateway handles all CORS
```

**Why:**
- Prevents duplicate headers
- Single source of truth
- Easier to manage origins

### 4. Different Middleware for Dev/Prod

**Development:**
```go
// Backend validates JWT locally
http.HandleFunc("/api/posts", middleware.DevAuthMiddleware(handlers.HandlePosts))
```

**Production:**
```go
// Backend reads X-User-* headers from Gateway
http.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.HandlePosts))
```

### 5. Same JWT_SECRET Everywhere

**CRITICAL:** JWT_SECRET must match in Gateway and Backend

```bash
# Check they match
diff <(cat gateway/.env | grep JWT_SECRET) \
     <(cat backend/.env | grep JWT_SECRET)
# Should show no difference
```

## Common Mistakes to Avoid

### ❌ Don't Do This

1. **Direct fetch with hardcoded URL:**
   ```typescript
   const response = await fetch('http://localhost:8000/api/posts');
   ```

2. **Backend sets CORS headers:**
   ```go
   w.Header().Set("Access-Control-Allow-Origin", "*")
   ```

3. **Forget credentials:**
   ```typescript
   fetch('/api/posts')  // Missing credentials: 'include'
   ```

4. **Different JWT_SECRET:**
   ```bash
   # Gateway
   JWT_SECRET=secret1
   
   # Backend
   JWT_SECRET=secret2  # ❌ Different!
   ```

5. **Use DevAuthMiddleware in production:**
   ```go
   // ❌ Wrong for production
   http.HandleFunc("/api/posts", middleware.DevAuthMiddleware(handlers.HandlePosts))
   ```

### ✅ Do This Instead

1. **Use apiClient:**
   ```typescript
   import { apiClient } from '@/lib/api';
   const response = await apiClient.get('/api/posts');
   ```

2. **Let Gateway manage CORS:**
   ```go
   // Backend doesn't set CORS at all
   ```

3. **apiClient handles credentials:**
   ```typescript
   // apiClient automatically adds credentials: 'include'
   const response = await apiClient.get('/api/posts');
   ```

4. **Same JWT_SECRET:**
   ```bash
   # Copy from Gateway to Backend
   JWT_SECRET=$(cat gateway/.env | grep JWT_SECRET | cut -d'=' -f2)
   echo "JWT_SECRET=$JWT_SECRET" >> backend/.env
   ```

5. **Use correct middleware:**
   ```go
   // Development
   http.HandleFunc("/api/posts", middleware.DevAuthMiddleware(handlers.HandlePosts))
   
   // Production
   http.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.HandlePosts))
   ```

## Testing

### Quick Test Commands

```bash
# 1. Test Gateway health
curl http://localhost/health

# 2. Test Backend health
curl http://localhost:8000/api/health

# 3. Test login
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}' \
  -c cookies.txt

# 4. Test authenticated request
curl http://localhost/api/posts \
  -b cookies.txt

# 5. Test CORS preflight
curl -X OPTIONS http://localhost/api/posts \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

### Verification Checklist

- [ ] Login works (200 OK)
- [ ] Cookie is set after login
- [ ] Protected routes work (not 401)
- [ ] No CORS errors in console
- [ ] Preflight returns 200 OK
- [ ] Access-Control-Allow-Origin NOT duplicated
- [ ] JWT_SECRET matches everywhere

## Maintenance

### Regular Tasks

**Weekly:**
- Check logs for authentication errors
- Check logs for CORS errors
- Verify all services healthy

**Monthly:**
- Review and update documentation
- Check for security updates
- Review JWT_SECRET rotation policy

**When Adding New Service:**
1. Add service URL to Gateway .env
2. Add routing in Gateway
3. Use AuthMiddleware in service
4. Don't set CORS in service
5. Test authentication works

### Adding New Frontend Origin

```go
// gateway/middleware/cors.go

allowedOrigins := map[string]bool{
    "http://localhost:3000": true,
    "https://new-frontend.com": true,  // ← Add here
}
```

## Support

### Getting Help

1. **Check Troubleshooting Guide:**
   - Common errors and solutions
   - Diagnostic commands
   - Verification steps

2. **Check Logs:**
   ```bash
   # Backend logs
   tail -f /tmp/main-backend.log
   
   # Gateway logs
   tail -f gateway/logs/gateway.log
   ```

3. **Collect Diagnostic Info:**
   ```bash
   # Service status
   curl http://localhost/health
   curl http://localhost:8000/api/health
   
   # Environment
   cat frontend/.env.local
   cat backend/.env | grep JWT_SECRET
   
   # CORS test
   curl -X OPTIONS http://localhost/api/posts \
     -H "Origin: http://localhost:3000" -v
   ```

### Bug Report Template

```markdown
## Error Description
[Describe the error]

## Environment
- [ ] Development
- [ ] Production

## Steps to Reproduce
1. ...
2. ...
3. ...

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happens]

## Logs
```
[Paste relevant logs]
```

## Screenshots
[Attach screenshots if applicable]
```

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 03.02.2026 | Initial specification created |

## Related Documentation

- [DEPLOYMENT.md](../../../DEPLOYMENT.md) - Deployment instructions
- [gateway.md](../../../gateway.md) - Gateway documentation
- [workspace-rules.md](../../steering/workspace-rules.md) - Workspace rules
- [auth-flow.md](../../steering/auth-flow.md) - Authentication flow rules

---

**For questions or updates, refer to the individual documents or contact the development team.**
