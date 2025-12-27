-- Добавляем поля для запланированных постов
ALTER TABLE posts ADD COLUMN status TEXT DEFAULT 'published' CHECK(status IN ('published', 'scheduled', 'draft'));
ALTER TABLE posts ADD COLUMN scheduled_at TEXT;

-- Обновляем существующие посты
UPDATE posts SET status = 'published' WHERE status IS NULL;
