-- Таблица для хранения медиа-файлов пользователей
CREATE TABLE IF NOT EXISTS user_media (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  file_name TEXT NOT NULL,           -- UUID имя файла (uuid.jpg)
  original_name TEXT NOT NULL,       -- Оригинальное имя файла
  file_path TEXT NOT NULL,           -- Относительный путь (users/1/photos/2025/12/uuid.jpg)
  file_size INTEGER NOT NULL,        -- Размер в байтах
  mime_type TEXT NOT NULL,           -- image/jpeg, video/mp4, etc.
  media_type TEXT NOT NULL,          -- 'photo', 'video', 'document', 'avatar'
  width INTEGER,                     -- Ширина изображения
  height INTEGER,                    -- Высота изображения
  duration INTEGER,                  -- Длительность видео (секунды)
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_media_user_id ON user_media(user_id);
CREATE INDEX IF NOT EXISTS idx_user_media_type ON user_media(media_type);
CREATE INDEX IF NOT EXISTS idx_user_media_uploaded ON user_media(uploaded_at);
