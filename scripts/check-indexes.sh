#!/bin/bash

# Скрипт для проверки использования индексов в запросах
# Использует EXPLAIN QUERY PLAN для анализа

DB_PATH="database/data.db"

echo "🔍 Проверка использования индексов в критичных запросах"
echo "=========================================================="
echo ""

# Функция для проверки запроса
check_query() {
    local name=$1
    local query=$2
    
    echo "📊 $name"
    echo "---"
    sqlite3 "$DB_PATH" "EXPLAIN QUERY PLAN $query" | while read line; do
        if [[ $line == *"SCAN TABLE"* ]]; then
            echo "❌ $line (полное сканирование таблицы - плохо!)"
        elif [[ $line == *"SEARCH TABLE"* ]] && [[ $line == *"USING INDEX"* ]]; then
            echo "✅ $line (использует индекс - хорошо!)"
        else
            echo "   $line"
        fi
    done
    echo ""
}

# 1. Загрузка ленты постов
check_query "Лента постов (главная страница)" \
"SELECT * FROM posts 
WHERE is_deleted = 0 AND status = 'published' 
ORDER BY created_at DESC 
LIMIT 20;"

# 2. Комментарии к посту
check_query "Комментарии к посту" \
"SELECT * FROM comments 
WHERE post_id = 1 
ORDER BY created_at DESC;"

# 3. Проверка лайка
check_query "Проверка лайка пользователя" \
"SELECT * FROM likes 
WHERE user_id = 1 AND post_id = 1;"

# 4. Подсчёт лайков
check_query "Подсчёт лайков поста" \
"SELECT COUNT(*) FROM likes 
WHERE post_id = 1;"

# 5. Непрочитанные уведомления
check_query "Непрочитанные уведомления" \
"SELECT * FROM notifications 
WHERE user_id = 1 AND is_read = 0 
ORDER BY created_at DESC;"

# 6. Роли пользователя
check_query "Активные роли пользователя" \
"SELECT * FROM user_roles 
WHERE user_id = 1 AND is_active = 1;"

# 7. Организации по типу и региону
check_query "Поиск организаций по типу и региону" \
"SELECT * FROM organizations 
WHERE type = 'shelter' AND address_region = 'Москва';"

# 8. Участники организации
check_query "Участники организации" \
"SELECT * FROM organization_members 
WHERE organization_id = 1;"

# 9. Посты с питомцами
check_query "Посты с конкретным питомцем" \
"SELECT p.* FROM posts p 
INNER JOIN post_pets pp ON p.id = pp.post_id 
WHERE pp.pet_id = 1;"

# 10. Друзья пользователя
check_query "Друзья пользователя" \
"SELECT * FROM friendships 
WHERE (user_id = 1 OR friend_id = 1) AND status = 'accepted';"

echo "=========================================================="
echo "✅ Проверка завершена!"
echo ""
echo "Легенда:"
echo "  ✅ SEARCH TABLE USING INDEX - запрос использует индекс (быстро)"
echo "  ❌ SCAN TABLE - полное сканирование таблицы (медленно)"
echo ""
echo "Если видите SCAN TABLE - нужно добавить индекс!"
