#!/bin/bash

# ============================================
# Скрипт восстановления базы данных из резервной копии
# ============================================

set -e  # Остановка при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
DB_PATH="database/data.db"
BACKUP_DIR="database/backups"

echo -e "${BLUE}🔄 Восстановление базы данных из резервной копии${NC}"
echo ""

# Проверка существования папки с бэкапами
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}❌ Ошибка: Папка с резервными копиями не найдена: $BACKUP_DIR${NC}"
    exit 1
fi

# Список доступных резервных копий
echo -e "${YELLOW}📁 Доступные резервные копии (.db файлы):${NC}"
echo ""

DB_BACKUPS=($(ls -1t "$BACKUP_DIR"/data_*.db 2>/dev/null))

if [ ${#DB_BACKUPS[@]} -eq 0 ]; then
    echo -e "${RED}❌ Резервные копии не найдены${NC}"
    exit 1
fi

# Показываем список с номерами
for i in "${!DB_BACKUPS[@]}"; do
    FILE="${DB_BACKUPS[$i]}"
    SIZE=$(du -h "$FILE" | cut -f1)
    DATE=$(basename "$FILE" | sed 's/data_\(.*\)\.db/\1/')
    # Форматируем дату для читаемости
    FORMATTED_DATE=$(echo "$DATE" | sed 's/\([0-9]\{4\}\)\([0-9]\{2\}\)\([0-9]\{2\}\)_\([0-9]\{2\}\)\([0-9]\{2\}\)\([0-9]\{2\}\)/\1-\2-\3 \4:\5:\6/')
    echo -e "  ${GREEN}[$((i+1))]${NC} $FORMATTED_DATE (размер: $SIZE)"
done

echo ""
echo -e "${YELLOW}📁 Доступные SQL dumps:${NC}"
echo ""

SQL_BACKUPS=($(ls -1t "$BACKUP_DIR"/backup_*.sql 2>/dev/null))

if [ ${#SQL_BACKUPS[@]} -eq 0 ]; then
    echo -e "${YELLOW}⚠️  SQL dumps не найдены${NC}"
else
    for i in "${!SQL_BACKUPS[@]}"; do
        FILE="${SQL_BACKUPS[$i]}"
        SIZE=$(du -h "$FILE" | cut -f1)
        DATE=$(basename "$FILE" | sed 's/backup_\(.*\)\.sql/\1/')
        FORMATTED_DATE=$(echo "$DATE" | sed 's/\([0-9]\{4\}\)\([0-9]\{2\}\)\([0-9]\{2\}\)_\([0-9]\{2\}\)\([0-9]\{2\}\)\([0-9]\{2\}\)/\1-\2-\3 \4:\5:\6/')
        echo -e "  ${GREEN}[S$((i+1))]${NC} $FORMATTED_DATE (размер: $SIZE)"
    done
fi

echo ""
echo -e "${YELLOW}Выберите резервную копию для восстановления:${NC}"
echo -e "  - Введите номер [1-${#DB_BACKUPS[@]}] для .db файла"
if [ ${#SQL_BACKUPS[@]} -gt 0 ]; then
    echo -e "  - Введите S1-S${#SQL_BACKUPS[@]} для SQL dump"
fi
echo -e "  - Введите 'q' для отмены"
echo ""
read -p "Ваш выбор: " CHOICE

# Обработка выбора
if [ "$CHOICE" = "q" ] || [ "$CHOICE" = "Q" ]; then
    echo -e "${YELLOW}❌ Восстановление отменено${NC}"
    exit 0
fi

# Проверка выбора SQL dump
if [[ "$CHOICE" =~ ^[Ss][0-9]+$ ]]; then
    INDEX=$(echo "$CHOICE" | sed 's/[Ss]//')
    if [ "$INDEX" -lt 1 ] || [ "$INDEX" -gt ${#SQL_BACKUPS[@]} ]; then
        echo -e "${RED}❌ Неверный номер резервной копии${NC}"
        exit 1
    fi
    
    BACKUP_FILE="${SQL_BACKUPS[$((INDEX-1))]}"
    RESTORE_TYPE="sql"
    
elif [[ "$CHOICE" =~ ^[0-9]+$ ]]; then
    if [ "$CHOICE" -lt 1 ] || [ "$CHOICE" -gt ${#DB_BACKUPS[@]} ]; then
        echo -e "${RED}❌ Неверный номер резервной копии${NC}"
        exit 1
    fi
    
    BACKUP_FILE="${DB_BACKUPS[$((CHOICE-1))]}"
    RESTORE_TYPE="db"
else
    echo -e "${RED}❌ Неверный выбор${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}⚠️  ВНИМАНИЕ!${NC}"
echo -e "Текущая база данных будет заменена на резервную копию:"
echo -e "  ${BLUE}$BACKUP_FILE${NC}"
echo ""
read -p "Вы уверены? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}❌ Восстановление отменено${NC}"
    exit 0
fi

# Создаём резервную копию текущей БД перед восстановлением
echo ""
echo -e "${YELLOW}📦 Создание резервной копии текущей БД перед восстановлением...${NC}"
if [ -f "$DB_PATH" ]; then
    SAFETY_BACKUP="${BACKUP_DIR}/before_restore_$(date +%Y%m%d_%H%M%S).db"
    cp "$DB_PATH" "$SAFETY_BACKUP"
    echo -e "${GREEN}✅ Текущая БД сохранена: $SAFETY_BACKUP${NC}"
fi

# Восстановление
echo ""
if [ "$RESTORE_TYPE" = "db" ]; then
    echo -e "${YELLOW}🔄 Восстановление из .db файла...${NC}"
    cp "$BACKUP_FILE" "$DB_PATH"
else
    echo -e "${YELLOW}🔄 Восстановление из SQL dump...${NC}"
    # Удаляем текущую БД
    rm -f "$DB_PATH"
    # Восстанавливаем из SQL
    sqlite3 "$DB_PATH" < "$BACKUP_FILE"
fi

# Проверка целостности
echo -e "${YELLOW}🔍 Проверка целостности восстановленной БД...${NC}"
INTEGRITY=$(sqlite3 "$DB_PATH" "PRAGMA integrity_check;" 2>&1)

if [ "$INTEGRITY" = "ok" ]; then
    echo -e "${GREEN}✅ База данных прошла проверку целостности${NC}"
else
    echo -e "${RED}❌ Ошибка целостности: $INTEGRITY${NC}"
    echo -e "${YELLOW}⚠️  Восстанавливаем из safety backup...${NC}"
    cp "$SAFETY_BACKUP" "$DB_PATH"
    exit 1
fi

# Статистика восстановленной БД
echo ""
echo -e "${YELLOW}📊 Статистика восстановленной БД:${NC}"
sqlite3 "$DB_PATH" "
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

echo ""
echo -e "${GREEN}✅ База данных успешно восстановлена!${NC}"
echo ""
echo -e "${YELLOW}💡 Резервная копия текущей БД сохранена в:${NC}"
echo "   $SAFETY_BACKUP"
echo ""
