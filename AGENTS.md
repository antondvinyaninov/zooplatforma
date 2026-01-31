# AGENTS.md

This repository is a multi-service monorepo (Next.js + Go + shared TS + Expo). Use this guide when acting as an agentic coding assistant.

## 1) Build / Lint / Test Commands

### Global (recommended)
- `make help` — list all repo commands.
- `make install` — install all frontend deps for all services.
- `make install-go` — download all Go deps for all services.
- `make install-all` — install npm + Go deps.
- `make all` — start all services (same as `./run`).
- `make stop` — stop all dev services.
- `make lint-go` — run golangci-lint on main backend (warning if missing).
- `make lint-ts` — run ESLint on main frontend.
- `make format-go` — gofmt across repo.
- `make format-ts` — runs `npm run format` in main frontend (if defined).

### Per-service dev/build
- Frontend (example: `main/frontend`):
  - `npm run dev` — Next.js dev server.
  - `npm run build` — Next.js production build.
  - `npm run lint` — ESLint.
- Backend (example: `main/backend`):
  - `air` — hot reload (uses `.air.toml`).
  - `go run main.go` — start server directly.

### Docker
- `make docker-build` — build all images.
- `make docker-up` — run all containers.
- `make docker-down` — stop all containers.
- `make docker-logs` — follow logs.

### Tests
- `make test` — runs performance test script.
- `make test-api` — API test (`tests/api-test.sh`).
- `make test-sso` — SSO integration test (`tests/sso-test.sh`).
- `make test-performance` — performance test.
- `./tests/performance-test.sh` — direct script.
- `./tests/compare-with-vk.sh` — comparison benchmark.
- `./tests/web-vitals-test.sh` — Lighthouse Web Vitals.

### Running a single test
- Shell scripts are the primary test units here:
  - `./tests/sso-test.sh` (SSO only)
  - `./tests/performance-test.sh` (perf only)
  - `./tests/api-test.sh` (API only)
- Planned unit tests (not always present):
  - Go: `go test ./...` within a backend folder.
  - Frontend: `npm test` in a frontend folder (if configured).

## 2) Code Style & Architecture Guidelines

### Project structure and paths
- Each service has its own `backend/` (Go) and `frontend/` (Next.js).
- Always verify service + layer before editing (example: `clinic/backend/...` vs `main/backend/...`).
- Ports (backend/frontend): main 8000/3000, admin 9000/4000, petbase 8100/4100,
  shelter 8200/5100, owner 8400/6100, volunteer 8500/6200, clinic 8600/6300.

### Service ownership & data rules
- Single shared SQLite DB at `database/data.db` (used by all services).
- PetBase is the source of truth for animal data; other services should not mutate
  PetBase tables directly.
- Reading other services’ tables is OK (direct SQL); mutating other services’ owned
  data must be via their REST API.

### API design
- REST style: `/api/<resource>` with plural nouns and kebab-case.
- Standard response:
  - Success: `{ "success": true, "data": ... }`
  - Error: `{ "success": false, "error": "...", "details"?: {...} }`
- Use correct HTTP codes (200/201/400/401/403/404/409/500).

### Error handling (backend Go)
- Never ignore errors; log them with context.
- Return user-friendly messages without internal details.
- Use consistent patterns and correct status codes.
- Prefer helper functions for consistent JSON responses.

### Error handling (frontend TS)
- Wrap API calls in `try/catch`.
- Surface errors in UI (inline/Toast).
- Handle 401/403 with redirects.
- Consider Error Boundaries for React rendering issues.

### SSO & Auth
- All services use SSO via Main backend.
- JWT token retrieved from Main (`/api/auth/me`).
- Frontends must use `credentials: 'include'` and pass Bearer token when needed.
- Backends must implement `AuthMiddleware` and set context keys.

### Shared components (`shared/`)
- `shared/` contains types, hooks, utils, API clients, and components.
- Tailwind v4 doesn’t scan node_modules by default; ensure each frontend adds
  shared paths to `tailwind.config.ts` content. If styles still fail, copy the
  component locally.
- Follow `shared/README.md` and `.kiro/steering/shared-components.md`.

### Imports & module usage
- Keep imports explicit; avoid deep relative chains if a shared module exists.
- Use `@pet/shared` for shared types/utils/API where allowed.
- For UI components with Tailwind styles, verify Tailwind config or copy locally.

### Formatting and types
- Go: standard `gofmt`; idiomatic error checks (`if err != nil { ... }`).
- TypeScript: prefer strict typing and explicit interfaces for API responses.
- Avoid `any` unless truly necessary and documented.

### Naming conventions
- Endpoints: plural nouns, kebab-case.
- DB tables: follow existing schema conventions.
- File names: lower-case with hyphens or Go standard file names.

### Logging
- Use consistent log formatting with context (IDs, endpoint, service).
- Emojis are used in some logging rules for quick scanning.

### Tests & scripts
- Keep test scripts executable (`chmod +x`).
- Update tests documentation if you add new scripts.

## 3) Cursor / Copilot rules
- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` were found in this repo.

## 4) References (key docs)
- `README.md` — full overview & quick start
- `Makefile` — canonical commands
- `tests/README.md` — test scripts & requirements
- `infrastructure/README.md` — docker usage
- `.kiro/steering/*` — architecture & style rules
- `shared/README.md` — shared package usage
