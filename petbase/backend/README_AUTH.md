# Аутентификация в ЗооБазе

## Для разработки

### Вариант 1: Использование X-User-ID заголовка (самый простой)

```bash
# Получить всех питомцев пользователя 1
curl -H "X-User-ID: 1" http://localhost:8100/api/pets

# Создать питомца
curl -X POST \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{"name":"Барсик","species":"Кошка"}' \
  http://localhost:8100/api/pets

# Обновить питомца (только своего!)
curl -X PUT \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{"name":"Барсик Обновлённый","species":"Кошка"}' \
  http://localhost:8100/api/pets/1
```

### Вариант 2: Использование JWT токена

```javascript
// 1. Сгенерировать токен (в Go коде или через отдельный endpoint)
import { middleware } from 'petbase/middleware';

const token = middleware.GenerateToken(1); // user_id = 1

// 2. Использовать токен в запросах
fetch('http://localhost:8100/api/pets', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

### Вариант 3: Через API клиент (рекомендуется)

```typescript
// main/frontend/lib/petbase-api.ts

import { petBaseAPI } from '@/lib/petbase-api';

// Установить user_id для разработки
petBaseAPI.setUserId(1);

// Или установить JWT токен
petBaseAPI.setAuthToken('your-jwt-token');

// Теперь все запросы будут аутентифицированы
const pets = await petBaseAPI.getPets();
const pet = await petBaseAPI.createPet({ name: 'Барсик', species: 'Кошка' });
```

## Для продакшена

### 1. Настройка JWT Secret

```bash
# Создайте .env файл
cp .env.example .env

# Сгенерируйте сильный секретный ключ (минимум 32 символа)
openssl rand -base64 32

# Добавьте в .env
JWT_SECRET=ваш-сгенерированный-ключ
```

### 2. Настройка CORS

```bash
# В .env укажите разрешённые origins
ALLOWED_ORIGINS=https://yoursite.com,https://petbase.yoursite.com
```

### 3. Интеграция с основным сайтом

```typescript
// После успешного логина на основном сайте
const response = await fetch('/api/auth/login', {
  method: 'POST',
  body: JSON.stringify({ email, password })
});

const { token } = await response.json();

// Сохранить токен
petBaseAPI.setAuthToken(token);

// Теперь все запросы к ЗооБазе будут аутентифицированы
```

## Проверка безопасности

### Тест 1: Попытка редактировать чужого питомца

```bash
# Создать питомца от пользователя 1
curl -X POST \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{"name":"Питомец пользователя 1","species":"Собака"}' \
  http://localhost:8100/api/pets

# Попытаться изменить от пользователя 2 (должно вернуть 403 Forbidden)
curl -X PUT \
  -H "X-User-ID: 2" \
  -H "Content-Type: application/json" \
  -d '{"name":"Взломанное имя","species":"Собака"}' \
  http://localhost:8100/api/pets/1
```

### Тест 2: Запрос без аутентификации

```bash
# Должно вернуть 401 Unauthorized
curl http://localhost:8100/api/pets
```

### Тест 3: Запрос с неразрешённого origin

```bash
# Должно вернуть 403 Forbidden
curl -H "Origin: http://evil-site.com" http://localhost:8100/api/pets
```

## Endpoints

### Публичные (без аутентификации)
- `GET /api/health` - проверка здоровья
- `GET /api/species` - список видов
- `GET /api/breeds` - список пород
- `GET /api/cards` - информационные карточки

### Защищённые (требуют аутентификации)
- `GET /api/pets` - список питомцев пользователя
- `POST /api/pets` - создать питомца
- `GET /api/pets/:id` - получить питомца
- `PUT /api/pets/:id` - обновить питомца (только своего!)
- `DELETE /api/pets/:id` - удалить питомца (только своего!)
- `GET /api/pets/user/:userId` - питомцы пользователя
- `GET /api/pets/status/:status` - питомцы по статусу

## Коды ошибок

- `401 Unauthorized` - не предоставлен токен или токен невалидный
- `403 Forbidden` - попытка доступа к чужому ресурсу или неразрешённый origin
- `404 Not Found` - ресурс не найден
- `400 Bad Request` - невалидные данные

## Логирование

Все заблокированные запросы логируются:

```
⚠️ Blocked request from unauthorized origin: http://evil-site.com
```

Проверяйте логи сервера для мониторинга подозрительной активности.
