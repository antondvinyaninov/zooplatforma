# Database Seeds

Скрипты для заполнения базы данных тестовыми данными.

## Структура

```
seeds/
├── dev/              # Данные для разработки
│   ├── users.sql
│   ├── posts.sql
│   └── comments.sql
├── test/             # Данные для тестов
│   └── test-data.sql
└── README.md
```

## Использование

### Для разработки

```bash
./seed-dev.sh
```

Заполняет БД тестовыми пользователями, постами и комментариями.

### Для тестов

```bash
./seed-test.sh
```

Заполняет БД минимальным набором данных для тестов.

## Пример seed файла

**users.sql:**
```sql
INSERT INTO users (name, email, password_hash, role) VALUES
('Тестовый Пользователь', 'test@example.com', '$2a$10$...', 'user'),
('Администратор', 'admin@example.com', '$2a$10$...', 'admin');
```

## Очистка данных

```bash
./clean-db.sh
```

Удаляет все данные из БД (кроме схемы).
