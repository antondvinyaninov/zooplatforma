---
inclusion: always
---

# Правила миграций базы данных

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Сохранение данных

**При миграции базы данных ВСЕГДА сохраняйте существующие данные!**

### Основные принципы

1. **Добавление полей** - старые данные остаются нетронутыми
2. **Изменение структуры** - данные мигрируют, но не удаляются
3. **Удаление полей** - только если поле действительно больше не нужно

### ✅ Правильные миграции

```sql
-- Добавление нового поля
ALTER TABLE users ADD COLUMN last_name TEXT DEFAULT '';

-- Разделение поля на два (данные сохраняются)
UPDATE users 
SET last_name = CASE 
    WHEN instr(name, ' ') > 0 THEN substr(name, instr(name, ' ') + 1)
    ELSE ''
END,
name = CASE 
    WHEN instr(name, ' ') > 0 THEN substr(name, 1, instr(name, ' ') - 1)
    ELSE name
END
WHERE name LIKE '% %';
-- ✅ Поле name сохраняется, добавляется last_name

-- Переименование поля (с сохранением данных)
ALTER TABLE users RENAME COLUMN old_field TO new_field;

-- Изменение типа данных (с конвертацией)
ALTER TABLE posts ADD COLUMN views_count INTEGER DEFAULT 0;
UPDATE posts SET views_count = CAST(views AS INTEGER);
-- После проверки можно удалить старое поле
```

### ❌ Неправильные миграции

```sql
-- НЕ ДЕЛАЙТЕ ТАК: затирание существующих данных
UPDATE users SET avatar = '';
-- ❌ Удаляет все аватарки!

-- НЕ ДЕЛАЙТЕ ТАК: удаление поля без миграции данных
ALTER TABLE users DROP COLUMN avatar;
-- ❌ Теряются все аватарки!

-- НЕ ДЕЛАЙТЕ ТАК: перезапись данных без проверки
UPDATE users SET name = 'Default';
-- ❌ Теряются все имена!
```

## Процесс создания миграции

### 1. Планирование

- Определите, какие данные затронуты
- Проверьте, нужно ли сохранить старые данные
- Продумайте обратную совместимость

### 2. Написание миграции

```sql
-- Шаг 1: Добавить новые поля
ALTER TABLE users ADD COLUMN new_field TEXT DEFAULT '';

-- Шаг 2: Мигрировать данные (если нужно)
UPDATE users SET new_field = old_field WHERE old_field IS NOT NULL;

-- Шаг 3: Удалить старые поля (ТОЛЬКО если они больше не нужны)
-- ALTER TABLE users DROP COLUMN old_field;
-- ⚠️ Закомментируйте удаление до полной проверки!
```

### 3. Тестирование

```bash
# Создайте резервную копию
cp database/data.db database/data.db.backup

# Примените миграцию
sqlite3 database/data.db < database/migrations/XXX_migration.sql

# Проверьте данные
sqlite3 database/data.db "SELECT * FROM users LIMIT 5;"

# Если что-то пошло не так, восстановите
# mv database/data.db.backup database/data.db
```

### 4. Проверка после миграции

- [ ] Все старые данные сохранены
- [ ] Новые поля заполнены корректно
- [ ] Приложение работает с новой структурой
- [ ] Нет ошибок в логах

## Типы миграций

### Добавление поля

```sql
-- Всегда безопасно
ALTER TABLE users ADD COLUMN phone TEXT DEFAULT '';
```

### Переименование поля

```sql
-- Безопасно, данные сохраняются
ALTER TABLE users RENAME COLUMN old_name TO new_name;
```

### Разделение поля

```sql
-- Добавляем новое поле
ALTER TABLE users ADD COLUMN last_name TEXT DEFAULT '';

-- Копируем данные (НЕ удаляем старое поле!)
UPDATE users 
SET last_name = substr(name, instr(name, ' ') + 1)
WHERE instr(name, ' ') > 0;

-- Обновляем старое поле (сохраняем первую часть)
UPDATE users 
SET name = substr(name, 1, instr(name, ' ') - 1)
WHERE instr(name, ' ') > 0;
```

### Объединение полей

```sql
-- Добавляем новое поле
ALTER TABLE users ADD COLUMN full_name TEXT DEFAULT '';

-- Объединяем данные
UPDATE users SET full_name = name || ' ' || last_name;

-- Старые поля НЕ удаляем до полной проверки!
```

### Удаление поля

```sql
-- ТОЛЬКО если поле действительно больше не используется!
-- Сначала убедитесь, что:
-- 1. Данные не нужны
-- 2. Код не использует это поле
-- 3. Есть резервная копия

-- ALTER TABLE users DROP COLUMN deprecated_field;
-- ⚠️ Всегда комментируйте удаление в первой версии миграции!
```

## Резервное копирование

### Перед каждой миграцией

```bash
# Создайте резервную копию с датой
DATE=$(date +%Y%m%d_%H%M%S)
cp database/data.db database/backups/data_${DATE}.db

# Или используйте sqlite3 dump
sqlite3 database/data.db .dump > database/backups/backup_${DATE}.sql
```

### Восстановление

```bash
# Из копии файла
cp database/backups/data_20241228_120000.db database/data.db

# Из SQL dump
sqlite3 database/data.db < database/backups/backup_20241228_120000.sql
```

## Checklist миграции

Перед применением миграции:

- [ ] Создана резервная копия базы данных
- [ ] Миграция не удаляет существующие данные
- [ ] Миграция не перезаписывает важные поля
- [ ] Добавлены DEFAULT значения для новых полей
- [ ] Протестирована на копии базы
- [ ] Проверена обратная совместимость
- [ ] Обновлены модели в коде (Go, TypeScript)
- [ ] Обновлены API handlers
- [ ] Обновлен frontend

После применения миграции:

- [ ] Проверены данные в базе
- [ ] Приложение запускается без ошибок
- [ ] Все функции работают корректно
- [ ] Нет потерянных данных
- [ ] Логи не содержат ошибок

## Примеры из проекта

### ✅ Хорошая миграция (006_add_user_media_table.sql)

```sql
-- Создание новой таблицы (не затрагивает существующие данные)
CREATE TABLE IF NOT EXISTS user_media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    file_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### ⚠️ Проблемная миграция (007_add_last_name_to_users.sql - первая версия)

```sql
-- Проблема: может затереть данные, если не аккуратно
UPDATE users 
SET name = CASE 
    WHEN instr(name, ' ') > 0 THEN substr(name, 1, instr(name, ' ') - 1)
    ELSE name
END;
-- ⚠️ Нужна дополнительная проверка, что данные не теряются
```

## Инструменты

### Проверка структуры таблицы

```bash
sqlite3 database/data.db "PRAGMA table_info(users);"
```

### Проверка данных

```bash
sqlite3 database/data.db "SELECT id, name, last_name, avatar FROM users LIMIT 10;"
```

### Подсчет записей

```bash
sqlite3 database/data.db "SELECT COUNT(*) FROM users WHERE avatar IS NOT NULL;"
```

---

**Помните:** Данные пользователей - это самое ценное в приложении. Всегда делайте резервные копии и тестируйте миграции!
