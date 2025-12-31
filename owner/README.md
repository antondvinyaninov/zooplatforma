# Owner - Кабинет владельца

Микросервис для управления питомцами владельцами.

## Порты

- **Frontend:** http://localhost:6100
- **Backend API:** http://localhost:8400

## Функционал

### Для владельцев:
- Список моих питомцев
- Просмотр информации о питомце
- История событий питомца
- Создание постов о питомцах

## Технологии

### Backend
- Go 1.25.5
- net/http (standard library)
- Air (hot reload)
- SQLite (shared database)
- JWT аутентификация через SSO

### Frontend
- Next.js 16.0.8 (App Router)
- React 19.2.1
- TypeScript 5.x
- Tailwind CSS 4.x

## Запуск

### Через главный скрипт (рекомендуется)
```bash
./run
```

### Отдельно Backend
```bash
cd owner/backend
air
```

### Отдельно Frontend
```bash
cd owner/frontend
npm install  # первый раз
npm run dev
```

## API Endpoints

### Защищённые (требуют JWT)

```
GET    /api/my-pets          Получить моих питомцев
GET    /api/pets/events      Получить события питомца
GET    /api/profile          Получить профиль пользователя
```

### Публичные

```
GET    /api/health           Health check
GET    /                     API info
```

## Авторизация

Используется SSO (Single Sign-On) через основной сайт:
- JWT токен в cookie `auth_token`
- Автоматическая проверка при каждом запросе
- Редирект на `/auth` если не авторизован

## Структура

```
owner/
├── backend/
│   ├── main.go              # Главный файл
│   ├── handlers/            # HTTP handlers
│   │   └── owner.go         # Owner endpoints
│   ├── middleware/          # Middleware
│   │   └── auth.go          # JWT аутентификация
│   ├── models/              # Модели данных
│   │   └── owner.go         # Pet, PetEvent, User
│   ├── .env                 # Переменные окружения
│   ├── .env.example         # Пример .env
│   ├── .air.toml            # Конфигурация Air
│   └── go.mod               # Go dependencies
└── frontend/
    ├── app/
    │   ├── page.tsx         # Главная (редирект)
    │   ├── auth/            # Авторизация
    │   ├── pets/            # Список питомцев
    │   └── layout.tsx       # Layout
    ├── lib/
    │   └── api.ts           # API клиент
    ├── package.json
    ├── tsconfig.json
    ├── tailwind.config.ts
    └── next.config.ts
```

## Переменные окружения

### Backend (.env)

```env
JWT_SECRET=<same_as_main_backend>
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:7000
```

**Важно:** JWT_SECRET должен быть таким же, как в main/backend/.env

## Интеграция с другими сервисами

### Main API (localhost:8000)
- Проверка авторизации (`/api/auth/me`)
- Получение информации о пользователе

### PetBase API (localhost:8100)
- Получение списка питомцев
- CRUD операции с питомцами

## Разработка

### Hot Reload

Backend автоматически перезагружается при изменении `.go` файлов (Air).
Frontend автоматически обновляется при изменении файлов (Next.js).

### Добавление нового endpoint

1. Добавьте handler в `backend/handlers/owner.go`
2. Зарегистрируйте route в `backend/main.go`
3. Добавьте метод в `frontend/lib/api.ts`
4. Используйте в компонентах

### Тестирование

```bash
# Проверка health
curl http://localhost:8400/api/health

# Получить питомцев (требует auth)
curl -H "Cookie: auth_token=<token>" \
  http://localhost:8400/api/my-pets
```

## Troubleshooting

### Backend не запускается

1. Проверьте, что порт 8400 свободен: `lsof -i :8400`
2. Проверьте `.env` файл (JWT_SECRET должен совпадать с main)
3. Проверьте логи: `cd owner/backend && air`

### Frontend не запускается

1. Проверьте, что порт 7000 свободен: `lsof -i :7000`
2. Установите зависимости: `cd owner/frontend && npm install`
3. Проверьте логи: `cd owner/frontend && npm run dev`

### Ошибка авторизации

1. Проверьте, что JWT_SECRET одинаковый в main и owner
2. Проверьте cookie `auth_token` в браузере
3. Попробуйте заново войти на основном сайте

## TODO

- [ ] Страница детальной карточки питомца `/pets/[id]`
- [ ] Создание постов о питомцах
- [ ] Редактирование информации о питомце
- [ ] История событий питомца
- [ ] Уведомления о событиях

## Документация

- [SSO Architecture](../docs/SSO_ARCHITECTURE.md)
- [API Endpoints](../docs/API_ENDPOINTS.md)
- [Main README](../README.md)
