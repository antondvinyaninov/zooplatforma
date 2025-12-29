# API Endpoints ЗооПлатформы

**Версия:** 1.0  
**Дата:** 29 декабря 2025  
**Статус:** Development

---

## 🌐 Микросервисы и порты

| Сервис | URL (Dev) | URL (Production) | Описание |
|--------|-----------|------------------|----------|
| **Main** | http://localhost:3000 | https://zooplatform.ru | Основной сайт (соцсеть) |
| **Main API** | http://localhost:8000 | https://api.zooplatform.ru | Backend основного сайта |
| **Admin** | http://localhost:4000 | https://admin.zooplatform.ru | Админ-панель |
| **Admin API** | http://localhost:9000 | https://admin-api.zooplatform.ru | Backend админ-панели |
| **PetID** | http://localhost:4100 | https://petid.zooplatform.ru | Реестр животных |
| **PetID API** | http://localhost:8100 | https://petid-api.zooplatform.ru | Backend реестра животных |
| **Shelter** | http://localhost:5000 | https://shelter.zooplatform.ru | Кабинеты организаций |
| **Shelter API** | http://localhost:8200 | https://shelter-api.zooplatform.ru | Backend организаций |
| **Mobile** | http://localhost:8081 | - | React Native (Expo) |

---

## 🔐 Авторизация

Все API используют **JWT токены** в cookie `auth_token`.

**SSO (Single Sign-On):** Вход на любом сервисе → токен доступен везде.

---

## 📱 Main API (localhost:8000)

### Авторизация

```
POST   /api/auth/register          Регистрация нового пользователя
POST   /api/auth/login             Вход в систему
POST   /api/auth/logout            Выход из системы
GET    /api/auth/me                Получить текущего пользователя
```

### Пользователи

```
GET    /api/users                  Получить всех пользователей
GET    /api/users/:id              Получить пользователя по ID
PUT    /api/profile                Обновить свой профиль
GET    /api/profile/:id            Получить профиль пользователя
```

### Посты

```
GET    /api/posts                  Получить все посты (лента)
GET    /api/posts/:id              Получить пост по ID
POST   /api/posts                  Создать пост
PUT    /api/posts/:id              Обновить свой пост
DELETE /api/posts/:id              Удалить свой пост
GET    /api/posts/user/:id         Получить посты пользователя
```

### Комментарии

```
GET    /api/posts/:id/comments     Получить комментарии к посту
POST   /api/posts/:id/comments     Добавить комментарий
DELETE /api/comments/:id           Удалить свой комментарий
```

### Лайки

```
POST   /api/posts/:id/like         Поставить лайк
DELETE /api/posts/:id/unlike       Убрать лайк
GET    /api/posts/:id/likes        Получить список лайкнувших
```

### Голосования (Polls)

```
POST   /api/polls                  Создать опрос
GET    /api/polls/:id              Получить опрос
POST   /api/polls/:id/vote         Проголосовать
GET    /api/polls/:id/results      Получить результаты
```

### Питомцы (базовое API, будет заменено на PetID)

```
GET    /api/pets/user/:id          Получить питомцев пользователя
POST   /api/pets                   Добавить питомца
DELETE /api/pets/:id               Удалить питомца
```

### Медиа (аватары, обложки, фото)

```
POST   /api/upload/avatar          Загрузить аватар
DELETE /api/upload/avatar          Удалить аватар
POST   /api/upload/cover           Загрузить обложку
DELETE /api/upload/cover           Удалить обложку
POST   /api/upload/media           Загрузить фото/видео
GET    /api/media/user/:id         Получить галерею пользователя
DELETE /api/media/:id              Удалить медиа
```

### Друзья

```
POST   /api/friends/request/:id    Отправить запрос в друзья
POST   /api/friends/accept/:id     Принять запрос
POST   /api/friends/reject/:id     Отклонить запрос
DELETE /api/friends/:id            Удалить из друзей
GET    /api/friends                Получить список друзей
GET    /api/friends/requests       Получить входящие запросы
```

### Подписки (планируется)

```
POST   /api/subscriptions/:id      Подписаться на пользователя
DELETE /api/subscriptions/:id      Отписаться
GET    /api/subscriptions          Мои подписки
GET    /api/subscribers            Мои подписчики
```

### Сообщения (планируется)

```
GET    /api/messages               Список диалогов
GET    /api/messages/:userId       История сообщений с пользователем
POST   /api/messages/:userId       Отправить сообщение
DELETE /api/messages/:id           Удалить сообщение
PUT    /api/messages/:id/read      Отметить как прочитанное
```

### Организации (DaData)

```
GET    /api/organizations/search   Поиск организаций по ИНН/ОГРН/названию
GET    /api/organizations/:inn     Получить информацию об организации
```

### Город по IP

```
GET    /api/location/city          Определить город по IP
```

### Аналитика (планируется)

```
GET    /api/analytics/user/:id     Статистика пользователя
```

### NPS опросы (планируется)

```
GET    /api/nps/should-show        Нужно ли показывать NPS опрос
POST   /api/nps/submit             Отправить ответ на NPS
POST   /api/nps/dismiss            Отложить опрос
POST   /api/nps/opt-out            Отписаться от опросов
```

---

## 👥 Admin API (localhost:9000)

### Пользователи

```
GET    /api/admin/users            Список всех пользователей
GET    /api/admin/users/:id        Информация о пользователе
PUT    /api/admin/users/:id/role   Изменить роль пользователя
DELETE /api/admin/users/:id        Удалить пользователя (бан)
POST   /api/admin/users/:id/unban  Разбанить пользователя
```

### Роли

```
GET    /api/admin/roles            Список ролей
POST   /api/admin/roles            Создать роль
PUT    /api/admin/roles/:id        Обновить роль
DELETE /api/admin/roles/:id        Удалить роль
```

### Модерация постов

```
GET    /api/admin/posts            Список всех постов
GET    /api/admin/posts/reported   Посты с жалобами
DELETE /api/admin/posts/:id        Удалить пост
POST   /api/admin/posts/:id/hide   Скрыть пост
POST   /api/admin/posts/:id/show   Показать пост
```

### Модерация сообщений

```
GET    /api/admin/messages/reported    Сообщения с жалобами
GET    /api/admin/messages/:id         Просмотр сообщения (только по жалобе)
DELETE /api/admin/messages/:id         Удалить сообщение
```

### Статистика

```
GET    /api/admin/stats/users          Статистика пользователей
GET    /api/admin/stats/posts          Статистика постов
GET    /api/admin/stats/activity       Статистика активности
GET    /api/admin/stats/nps            NPS метрики
```

### Логи

```
GET    /api/admin/logs                 Системные логи
GET    /api/admin/logs/errors          Логи ошибок
GET    /api/admin/logs/access          Логи доступа модераторов
```

### NPS аналитика

```
GET    /api/admin/nps/dashboard        Дашборд NPS
GET    /api/admin/nps/comments         Комментарии пользователей
GET    /api/admin/nps/export           Экспорт данных NPS
```

---

## 🐾 PetID API (localhost:8100)

### Питомцы (CRUD)

```
POST   /api/petid/pets                 Создать карточку питомца
GET    /api/petid/pets/:id             Получить полную карточку
GET    /api/petid/pets/:id/summary     Получить краткую информацию (для постов)
PUT    /api/petid/pets/:id             Обновить карточку
DELETE /api/petid/pets/:id             Удалить карточку
```

### Списки питомцев

```
GET    /api/petid/pets/user/:userId    Питомцы пользователя
GET    /api/petid/pets/search          Поиск питомцев
GET    /api/petid/pets/status/:status  Питомцы по статусу (home, looking_for_home, lost)
```

### События (история питомца)

```
GET    /api/petid/pets/:id/events      История событий питомца
POST   /api/petid/pets/:id/events      Добавить событие
PUT    /api/petid/events/:id           Обновить событие
DELETE /api/petid/events/:id           Удалить событие
```

**Типы событий:**
- `registration` - регистрация
- `ownership_change` - смена владельца
- `sterilization` - стерилизация
- `vaccination` - вакцинация
- `medical_visit` - визит к ветеринару
- `lost` - потерялся
- `found` - нашёлся
- `death` - смерть

### Чипы и метки

```
GET    /api/petid/chips/:number        Найти питомца по номеру чипа
POST   /api/petid/chips                Зарегистрировать чип
PUT    /api/petid/chips/:id            Обновить информацию о чипе
```

### Медицинская информация (планируется v1.1.0)

```
GET    /api/petid/pets/:id/medical     Медицинская карта
POST   /api/petid/pets/:id/medical     Добавить медицинскую запись
GET    /api/petid/pets/:id/vaccinations Список вакцинаций
POST   /api/petid/pets/:id/vaccinations Добавить вакцинацию
```

### Родословная (планируется v1.2.0)

```
GET    /api/petid/pets/:id/parents     Родители питомца
POST   /api/petid/pets/:id/parents     Указать родителей
GET    /api/petid/pets/:id/children    Потомки питомца
```

### Справочники (ЗооБаза)

```
GET    /api/petid/species              Список видов животных
GET    /api/petid/species/:id          Информация о виде
GET    /api/petid/breeds               Список пород
GET    /api/petid/breeds/:id           Информация о породе
GET    /api/petid/breeds/species/:id   Породы конкретного вида
```

---

## 🏥 Shelter API (localhost:8200) - планируется

### Приюты

```
GET    /api/shelter/list               Список приютов
GET    /api/shelter/:id                Информация о приюте
POST   /api/shelter                    Зарегистрировать приют
PUT    /api/shelter/:id                Обновить информацию
```

### Животные в приюте

```
GET    /api/shelter/:id/pets           Животные в приюте
POST   /api/shelter/:id/pets/intake    Принять животное
POST   /api/shelter/:id/pets/adopt     Пристроить животное
```

### Волонтёры

```
GET    /api/shelter/:id/volunteers     Волонтёры приюта
POST   /api/shelter/:id/volunteers     Добавить волонтёра
DELETE /api/shelter/:id/volunteers/:userId Удалить волонтёра
```

### Задачи

```
GET    /api/shelter/tasks              Мои задачи (волонтёр)
POST   /api/shelter/tasks              Создать задачу
PUT    /api/shelter/tasks/:id          Обновить задачу
DELETE /api/shelter/tasks/:id          Удалить задачу
```

### ОСВВ (Отлов-Стерилизация-Вакцинация-Возврат)

```
GET    /api/shelter/osvv/requests      Заявки на ОСВВ
POST   /api/shelter/osvv/requests      Подать заявку
PUT    /api/shelter/osvv/requests/:id  Обновить статус заявки
```

### Ветклиники

```
GET    /api/clinic/:id                 Информация о клинике
GET    /api/clinic/:id/schedule        Расписание клиники
POST   /api/clinic/:id/appointment     Записаться на приём
GET    /api/clinic/:id/appointments    Список записей
```

---

## 📊 Статические файлы

### Uploads (медиа пользователей)

```
GET    /uploads/avatars/:filename      Аватары пользователей
GET    /uploads/covers/:filename       Обложки профилей
GET    /uploads/posts/:filename        Медиа в постах
GET    /uploads/pets/:filename         Фото питомцев
```

**Development:** Локальная папка `./uploads/`  
**Production:** S3 / Yandex Object Storage + CDN

---

## 🔗 Внешние API

### DaData (поиск организаций)

```
Endpoint: https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party
Используется: Main API → /api/organizations/search
```

### Yandex Maps

```
Endpoint: https://api-maps.yandex.ru/2.1/
Используется: Frontend для отображения карт организаций
```

### Yandex Object Storage (планируется)

```
Endpoint: https://storage.yandexcloud.net
Bucket: zooplatform-uploads
Используется: Хранение медиа файлов в Production
```

---

## 🔄 WebSocket (планируется v1.0.0)

### Реал-тайм уведомления

```
WS     ws://localhost:8000/ws          WebSocket соединение
```

**События:**
- `new_message` - новое сообщение
- `new_comment` - новый комментарий
- `new_like` - новый лайк
- `friend_request` - запрос в друзья
- `pet_status_change` - изменение статуса питомца

---

## 📱 Mobile API

Mobile приложение использует те же API endpoints, что и веб-версия.

**Base URL (Dev):** http://localhost:8000  
**Base URL (Prod):** https://api.zooplatform.ru

---

## 🔐 Права доступа

### Роли пользователей:

| Роль | Описание | Доступ |
|------|----------|--------|
| `user` | Обычный пользователь | Main сайт |
| `volunteer` | Волонтёр | Main + Shelter (задачи) |
| `shelter_admin` | Администратор приюта | Main + Shelter (управление) |
| `clinic_admin` | Администратор клиники | Main + Clinic (управление) |
| `moderator` | Модератор | Main + Admin (модерация) |
| `superadmin` | Суперадминистратор | Всё |

### Проверка прав:

```go
// Middleware для проверки роли
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        if userRole != role && userRole != "superadmin" {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Использование
router.DELETE("/api/admin/posts/:id", RequireRole("moderator"), DeletePost)
```

---

## 📝 Формат ответов

### Успешный ответ:

```json
{
  "success": true,
  "data": {
    // данные
  }
}
```

### Ошибка:

```json
{
  "success": false,
  "error": "Описание ошибки"
}
```

### Пагинация:

```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

## 🚀 Rate Limiting

Для защиты от DDoS и злоупотреблений:

| Endpoint | Лимит |
|----------|-------|
| `/api/auth/login` | 5 запросов / минуту |
| `/api/auth/register` | 3 запроса / час |
| `/api/posts` (POST) | 10 постов / час |
| `/api/messages` (POST) | 100 сообщений / час |
| Остальные | 100 запросов / минуту |

---

## 📚 Документация API

**Swagger/OpenAPI:** (планируется)  
**Postman Collection:** (планируется)

---

## 🔮 Планируется в будущем

### v1.1.0 - Медицинская платформа
- `/api/petid/pets/:id/medical` - медицинские записи
- `/api/clinic/:id/appointments` - онлайн-запись к ветеринару

### v1.2.0 - Родословная
- `/api/petid/pets/:id/parents` - родители питомца
- `/api/petid/pets/:id/pedigree` - полная родословная

### v1.3.0 - ОСВВ и приюты
- `/api/shelter/osvv/*` - программа ОСВВ
- `/api/shelter/:id/visits` - онлайн-запись на посещение приюта

### v1.4.0 - Волонтёры и сборы
- `/api/fundraisers/*` - сборы средств
- `/api/volunteer/fosters` - передержки

### v2.0.0 - Расширенная функциональность
- GraphQL API
- WebSocket для real-time
- Elasticsearch для поиска

---

**Документ создан:** 29 декабря 2025  
**Автор:** ЗооПлатформа Team  
**Версия:** 1.0
