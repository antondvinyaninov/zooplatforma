# Full-Stack Project: Go + Next.js + React Native

Мультиплатформенное приложение с Go бекендом, Next.js веб-фронтендом и React Native мобильным приложением.

## Стек технологий

### Backend
- Go 1.25.5
- net/http (standard library)
- Air (hot reload)
- SQLite (database)

### Web Frontend
- Next.js 16.0.8 (App Router)
- React 19.2.1
- TypeScript 5.x
- Tailwind CSS 4.x

### Mobile
- React Native 0.81.5
- Expo 54.0.27
- TypeScript 5.x

### Shared
- Общая логика (API клиент, типы, хуки)
- TypeScript

## Структура проекта

```
├── backend/          # Go API сервер
│   ├── main.go
│   └── .air.toml     # Конфигурация hot reload
├── database/         # SQLite база данных (отдельный модуль)
│   ├── db.go         # Инициализация и работа с БД
│   └── data.db       # Файл базы данных
├── frontend/         # Next.js веб-приложение
│   └── src/
│       ├── components/
│       │   ├── ui/       # UI компоненты
│       │   ├── shared/   # Переиспользуемые компоненты
│       │   └── layout/   # Layout компоненты
│       └── lib/
│           └── api.ts    # API клиент
├── mobile/           # React Native приложение
│   ├── App.tsx
│   └── src/
│       └── lib/
│           └── api.ts    # API клиент
├── shared/           # Общая логика
│   └── src/
│       ├── api/      # API клиент
│       ├── types/    # TypeScript типы
│       ├── hooks/    # Переиспользуемые хуки
│       └── utils/    # Утилиты
└── run               # Скрипт запуска всего проекта
```

## Быстрый старт

```bash
./run
```

Это запустит все три сервиса:
- **Backend:** http://localhost:8080 (с hot reload)
- **Web:** http://localhost:3000
- **Mobile:** http://localhost:8081

Остановка: `Ctrl+C`

## API Endpoints

- `GET /` - Welcome message
- `GET /api/health` - Health check

## Разработка

### Backend (с hot reload)
```bash
cd backend
air
```

### Web Frontend
```bash
cd frontend
npm run dev
```

### Mobile
```bash
cd mobile
npm start
# Затем выбери платформу:
# w - web браузер
# i - iOS симулятор
# a - Android эмулятор
# или отсканируй QR код в Expo Go
```

## Особенности

- **Hot Reload:** Все три части проекта поддерживают автоматическую перезагрузку при изменениях
- **Shared Logic:** Общая бизнес-логика, типы и API клиент переиспользуются между платформами
- **CORS:** Настроен для работы со всеми клиентами
- **TypeScript:** Полная типизация во всех частях проекта
- **SQLite:** Легкая встроенная база данных, автоматически создается при первом запуске
- **Модульная архитектура:** База данных вынесена в отдельный модуль для переиспользования
