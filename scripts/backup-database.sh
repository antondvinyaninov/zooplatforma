#!/bin/bash

# ============================================
# Скрипт резервного копирования базы данных
# ============================================

set -e  # Остановка при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Конфигурация
DB_PATH="database/data.db"
BACKUP_DIR="database/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/data_${DATE}.db"
SQL_BACKUP="${BACKUP_DIR}/backup_${DATE}.sql"

# Максимальное количество резервных копий (старые удаляются)
MAX_BACKUPS=10

echo -e "${GREEN}🔄 Начинаем резервное копирование базы данных...${NC}"

# Проверка существования БД
if [ ! -f "$DB_PATH" ]; then
    echo -e "${RED}❌ Ошибка: База данных не найдена: $DB_PATH${NC}"
    exit 1
fi

# Создание папки для бэкапов если не существует
mkdir -p "$BACKUP_DIR"

# 1. Копирование файла БД
echo -e "${YELLOW}📦 Создание копии файла БД...${NC}"
cp "$DB_PATH" "$BACKUP_FILE"

if [ -f "$BACKUP_FILE" ]; then
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}✅ Файл БД скопирован: $BACKUP_FILE (размер: $SIZE)${NC}"
else
    echo -e "${RED}❌ Ошибка при копировании файла БД${NC}"
    exit 1
fi

# 2. Создание SQL dump
echo -e "${YELLOW}📝 Создание SQL dump...${NC}"
sqlite3 "$DB_PATH" .dump > "$SQL_BACKUP"

if [ -f "$SQL_BACKUP" ]; then
    SIZE=$(du -h "$SQL_BACKUP" | cut -f1)
    echo -e "${GREEN}✅ SQL dump создан: $SQL_BACKUP (размер: $SIZE)${NC}"
else
    echo -e "${RED}❌ Ошибка при создании SQL dump${NC}"
    exit 1
fi

# 3. Проверка целостности резервной копии
echo -e "${YELLOW}🔍 Проверка целостности резервной копии...${NC}"
INTEGRITY=$(sqlite3 "$BACKUP_FILE" "PRAGMA integrity_check;" 2>&1)

if [ "$INTEGRITY" = "ok" ]; then
    echo -e "${GREEN}✅ Резервная копия прошла проверку целостности${NC}"
else
    echo -e "${RED}❌ Ошибка целостности резервной копии: $INTEGRITY${NC}"
    exit 1
fi

# 4. Подсчёт записей в таблицах
echo -e "${YELLOW}📊 Статистика резервной копии:${NC}"
sqlite3 "$BACKUP_FILE" "
SELECT 
    'users: ' || COUNT(*) FROM users
UNION ALL
SELECT 
    'posts: ' || COUNT(*) FROM posts
UNION ALL
SELECT 
    'comments: ' || COUNT(*) FROM comments
UNION ALL
SELECT 
    'pets: ' || COUNT(*) FROM pets;
" 2>/dev/null || echo "Некоторые таблицы могут отсутствовать"

# 5. Удаление старых резервных копий
echo -e "${YELLOW}🗑️  Удаление старых резервных копий (оставляем последние $MAX_BACKUPS)...${NC}"

# Подсчёт количества .db файлов
DB_COUNT=$(ls -1 "$BACKUP_DIR"/data_*.db 2>/dev/null | wc -l | tr -d ' ')

if [ "$DB_COUNT" -gt "$MAX_BACKUPS" ]; then
    # Удаляем старые .db файлы
    ls -1t "$BACKUP_DIR"/data_*.db | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
    echo -e "${GREEN}✅ Удалено старых .db копий: $((DB_COUNT - MAX_BACKUPS))${NC}"
fi

# Подсчёт количества .sql файлов
SQL_COUNT=$(ls -1 "$BACKUP_DIR"/backup_*.sql 2>/dev/null | wc -l | tr -d ' ')

if [ "$SQL_COUNT" -gt "$MAX_BACKUPS" ]; then
    # Удаляем старые .sql файлы
    ls -1t "$BACKUP_DIR"/backup_*.sql | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
    echo -e "${GREEN}✅ Удалено старых .sql копий: $((SQL_COUNT - MAX_BACKUPS))${NC}"
fi

# 6. Итоговая информация
echo ""
echo -e "${GREEN}✅ Резервное копирование завершено успешно!${NC}"
echo ""
echo -e "${YELLOW}📁 Созданные файлы:${NC}"
echo "   - $BACKUP_FILE"
echo "   - $SQL_BACKUP"
echo ""
echo -e "${YELLOW}📊 Всего резервных копий в папке:${NC}"
echo "   - .db файлов: $(ls -1 "$BACKUP_DIR"/data_*.db 2>/dev/null | wc -l | tr -d ' ')"
echo "   - .sql файлов: $(ls -1 "$BACKUP_DIR"/backup_*.sql 2>/dev/null | wc -l | tr -d ' ')"
echo ""
echo -e "${YELLOW}💡 Для восстановления используйте:${NC}"
echo "   cp $BACKUP_FILE database/data.db"
echo "   или"
echo "   sqlite3 database/data.db < $SQL_BACKUP"
echo ""
