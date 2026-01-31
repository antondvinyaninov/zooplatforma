#!/bin/bash

# Скрипт миграции пользователей из database/data.db в auth/backend/auth.db
# Сохраняет ID пользователей для совместимости

set -e

echo "🔄 Миграция пользователей в Auth Service..."

# Проверка существования файлов
if [ ! -f "database/data.db" ]; then
    echo "❌ Файл database/data.db не найден"
    exit 1
fi

if [ ! -f "auth/backend/auth.db" ]; then
    echo "⚠️ Файл auth/backend/auth.db не найден, будет создан"
fi

# Получить список пользователей из основной БД
echo "📊 Получение пользователей из database/data.db..."
USERS=$(sqlite3 database/data.db "SELECT id, email, password, name FROM users;")

if [ -z "$USERS" ]; then
    echo "⚠️ Пользователи не найдены в database/data.db"
    exit 0
fi

# Подсчитать количество пользователей
USER_COUNT=$(echo "$USERS" | wc -l | tr -d ' ')
echo "✅ Найдено пользователей: $USER_COUNT"

# Мигрировать каждого пользователя
echo "🔄 Миграция пользователей..."

MIGRATED=0
SKIPPED=0

while IFS='|' read -r id email password name; do
    # Проверить существование пользователя в auth.db
    EXISTS=$(sqlite3 auth/backend/auth.db "SELECT COUNT(*) FROM users WHERE email = '$email';" 2>/dev/null || echo "0")
    
    if [ "$EXISTS" -gt 0 ]; then
        echo "⏭️  Пропущен: $email (уже существует)"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi
    
    # Получить аватар из старой БД
    AVATAR=$(sqlite3 database/data.db "SELECT avatar FROM users WHERE id = $id;" 2>/dev/null || echo "")
    
    # Вставить пользователя с сохранением ID и аватара
    sqlite3 auth/backend/auth.db "
        INSERT INTO users (id, email, password_hash, name, last_name, avatar, role, email_verified, created_at, updated_at)
        VALUES ($id, '$email', '$password', '$name', '', '$AVATAR', 'user', 0, datetime('now'), datetime('now'));
    " 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo "✅ Мигрирован: $email (ID: $id, Avatar: ${AVATAR:-нет})"
        MIGRATED=$((MIGRATED + 1))
    else
        echo "❌ Ошибка миграции: $email"
    fi
done <<< "$USERS"

echo ""
echo "📊 Итоги миграции:"
echo "   Всего пользователей: $USER_COUNT"
echo "   Мигрировано: $MIGRATED"
echo "   Пропущено: $SKIPPED"
echo ""
echo "✅ Миграция завершена!"
echo ""
echo "⚠️  ВАЖНО: Пароли скопированы как есть (уже захешированы)"
echo "   Пользователи могут войти с теми же паролями, что и раньше"
