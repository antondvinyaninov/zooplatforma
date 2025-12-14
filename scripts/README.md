# Scripts

Утилиты и скрипты для автоматизации задач.

## Структура

```
scripts/
├── migrations/       # Миграции базы данных
│   └── README.md
├── seeds/           # Тестовые данные
│   └── README.md
└── README.md        # Этот файл
```

## Миграции (migrations/)

Скрипты для изменения схемы базы данных.

### Создание миграции

```bash
cd scripts/migrations
./create-migration.sh "add_users_table"
```

### Применение миграций

```bash
cd scripts/migrations
./run-migrations.sh
```

## Seeds (seeds/)

Скрипты для заполнения БД тестовыми данными.

### Запуск seeds

```bash
cd scripts/seeds
./seed-dev.sh        # Для разработки
./seed-test.sh       # Для тестов
```

## Другие утилиты

- `backup-db.sh` - Резервное копирование БД
- `restore-db.sh` - Восстановление БД
- `clean-cache.sh` - Очистка кеша всех сервисов
- `check-ports.sh` - Проверка занятых портов
