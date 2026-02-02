# TODO: PostgreSQL Syntax Fixes

Файлы с SQL запросами использующими `?` плейсхолдеры, которые нужно конвертировать в `$1, $2, $3` для PostgreSQL.

## ✅ Исправлено:
- [x] `main/backend/handlers/users.go` - добавлен `ConvertPlaceholdersUsers()`
- [x] `main/backend/handlers/friends.go` - добавлен `convertPlaceholdersFriends()`
- [x] `main/backend/handlers/notifications.go` - добавлен `convertPlaceholdersNotif()`
- [x] `main/backend/handlers/organizations.go` - добавлен `convertPlaceholders()` в GetMyOrganizationsHandler
- [x] `main/backend/handlers/helpers.go` - создана глобальная функция `ConvertPlaceholders()`

## ⏳ Требуют исправления (по приоритету):

### Высокий приоритет (вызываются при загрузке главной страницы):
1. `main/backend/handlers/posts.go` - много запросов в getAllPosts, getUserPosts, createPost
2. `main/backend/handlers/likes.go` - toggleLike, getLikeStatus
3. `main/backend/handlers/comments.go` - создание/удаление комментариев
4. `main/backend/handlers/polls.go` - loadPollForPost, голосование

### Средний приоритет:
5. `main/backend/handlers/pets.go` - CRUD операции с питомцами
6. `main/backend/handlers/announcements.go` - объявления о питомцах
7. `main/backend/handlers/messenger.go` - сообщения
8. `main/backend/handlers/favorites.go` - избранное

### Низкий приоритет (редко используются):
9. `main/backend/handlers/user_activity.go` - статистика активности
10. `main/backend/handlers/user_logs.go` - логи пользователей

## Как исправлять:

Для каждого файла:
1. Добавить в начале файла: `import "backend/handlers"` (если нет)
2. Обернуть каждый SQL запрос с `?` в `ConvertPlaceholders()`:

```go
// Было:
db.Query("SELECT * FROM users WHERE id = ?", id)

// Стало:
query := ConvertPlaceholders("SELECT * FROM users WHERE id = ?")
db.Query(query, id)
```

3. Для `is_read = 1/0` заменить на `is_read = TRUE/FALSE`
4. Пересобрать: `go build -o test-build main.go` в `main/backend/`
5. Закоммитить и запушить

## Стратегия:
- Исправляем по мере появления ошибок в логах EasyPanel
- Начинаем с самых критичных (posts.go, likes.go)
- Используем глобальную функцию `ConvertPlaceholders()` из `helpers.go`
