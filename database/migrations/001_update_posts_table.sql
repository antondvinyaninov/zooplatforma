-- Миграция: Обновление таблицы posts для универсальной стены в стиле Threads
-- Дата: 2025-12-26

-- Шаг 1: Создаём новую таблицу с правильной структурой
CREATE TABLE IF NOT EXISTS posts_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  author_id INTEGER NOT NULL,
  author_type TEXT NOT NULL DEFAULT 'user', -- 'user' или 'organization'
  content TEXT DEFAULT '',
  attached_pets TEXT DEFAULT '[]', -- JSON массив PetID
  attachments TEXT DEFAULT '[]', -- JSON массив медиа-файлов
  tags TEXT DEFAULT '[]', -- JSON массив меток
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT 0,
  FOREIGN KEY (author_id) REFERENCES users(id)
);

-- Шаг 2: Копируем данные из старой таблицы
INSERT INTO posts_new (id, author_id, author_type, content, created_at, updated_at)
SELECT id, user_id, 'user', content, created_at, created_at
FROM posts;

-- Шаг 3: Удаляем старую таблицу
DROP TABLE posts;

-- Шаг 4: Переименовываем новую таблицу
ALTER TABLE posts_new RENAME TO posts;

-- Шаг 5: Создаём индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id, author_type);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_is_deleted ON posts(is_deleted);

-- Шаг 6: Создаём таблицу для связи постов и питомцев (для быстрых запросов)
CREATE TABLE IF NOT EXISTS post_pets (
  post_id INTEGER NOT NULL,
  pet_id INTEGER NOT NULL,
  PRIMARY KEY (post_id, pet_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_post_pets_pet ON post_pets(pet_id);
