-- Добавление поля anonymous_voting в таблицу polls

ALTER TABLE polls ADD COLUMN anonymous_voting BOOLEAN NOT NULL DEFAULT 0;
