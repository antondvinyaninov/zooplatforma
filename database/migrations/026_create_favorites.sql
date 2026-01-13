-- Миграция 026: Создание таблицы избранных питомцев
-- Дата: 2025-01-13
-- Описание: Таблица для хранения избранных питомцев пользователей

-- Таблица избранных питомцев
CREATE TABLE IF NOT EXISTS favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    UNIQUE(user_id, pet_id) -- Один питомец может быть добавлен в избранное только один раз
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_pet_id ON favorites(pet_id);
CREATE INDEX IF NOT EXISTS idx_favorites_created_at ON favorites(created_at);

-- Комментарии
-- user_id: ID пользователя, который добавил питомца в избранное
-- pet_id: ID питомца
-- created_at: Дата добавления в избранное
