-- Создаем таблицу лайков
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE(user_id, post_id) -- Один пользователь может лайкнуть пост только один раз
);

-- Индекс для быстрого поиска лайков поста
CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);

-- Индекс для быстрого поиска лайков пользователя
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
