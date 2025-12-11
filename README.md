# Next.js + Go Project

Full-stack приложение с Next.js фронтендом и Go бекендом.

## Стек технологий

### Backend
- Go 1.25.5
- net/http (standard library)

### Frontend
- Next.js 16.0.8
- React 19.2.1
- TypeScript 5.x
- Tailwind CSS 4.x

## Структура проекта

```
├── backend/          # Go API сервер
├── frontend/         # Next.js приложение
│   └── src/
│       ├── components/
│       │   ├── ui/       # UI компоненты
│       │   ├── shared/   # Переиспользуемые компоненты
│       │   └── layout/   # Layout компоненты
│       └── lib/
│           └── api.ts    # API клиент
└── run               # Скрипт запуска
```

## Быстрый старт

```bash
./run
```

Это запустит:
- Backend на http://localhost:8080
- Frontend на http://localhost:3000

Остановка: `Ctrl+C`

## API Endpoints

- `GET /` - Welcome message
- `GET /api/health` - Health check

## Разработка

### Backend
```bash
cd backend
go run main.go
```

### Frontend
```bash
cd frontend
npm run dev
```
