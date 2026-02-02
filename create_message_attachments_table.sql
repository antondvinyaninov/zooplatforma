-- Создание таблицы message_attachments для хранения вложений в сообщениях

CREATE TABLE IF NOT EXISTS message_attachments (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индекс для быстрого поиска вложений по message_id
CREATE INDEX IF NOT EXISTS idx_message_attachments_message_id ON message_attachments(message_id);
