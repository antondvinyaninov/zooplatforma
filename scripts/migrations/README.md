# Database Migrations

Миграции для управления схемой базы данных.

## Структура

Каждая миграция состоит из двух файлов:
- `YYYYMMDD_HHMMSS_description_up.sql` - применение изменений
- `YYYYMMDD_HHMMSS_description_down.sql` - откат изменений

## Создание миграции

```bash
./create-migration.sh "add_comments_table"
```

Это создаст:
- `20241214_120000_add_comments_table_up.sql`
- `20241214_120000_add_comments_table_down.sql`

## Применение миграций

```bash
./run-migrations.sh
```

## Откат миграции

```bash
./rollback-migration.sh
```

## Пример миграции

**up.sql:**
```sql
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**down.sql:**
```sql
DROP TABLE IF EXISTS comments;
```
