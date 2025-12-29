#!/bin/bash

# ============================================
# Скрипт резервного копирования файлов пользователей
# ============================================

set -e  # Остановка при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Конфигурация
UPLOADS_DIR="uploads"
BACKUP_DIR="database/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/uploads_${DATE}.tar.gz"

# Максимальное количество резервных копий
MAX_BACKUPS=5

echo -e "${GREEN}🔄 Начинаем резервное копирование файлов пользователей...${NC}"

# Проверка существования папки uploads
if [ ! -d "$UPLOADS_DIR" ]; then
    echo -e "${RED}❌ Ошибка: Папка uploads не найдена: $UPLOADS_DIR${NC}"
    exit 1
fi

# Создание папки для бэкапов если не существует
mkdir -p "$BACKUP_DIR"

# Подсчёт файлов
FILE_COUNT=$(find "$UPLOADS_DIR" -type f ! -name ".gitkeep" | wc -l | tr -d ' ')

if [ "$FILE_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Нет файлов для резервного копирования${NC}"
    exit 0
fi

echo -e "${YELLOW}📊 Найдено файлов: $FILE_COUNT${NC}"

# Подсчёт размера
SIZE=$(du -sh "$UPLOADS_DIR" | cut -f1)
echo -e "${YELLOW}📦 Размер: $SIZE${NC}"

# Создание архива
echo -e "${YELLOW}📦 Создание архива...${NC}"
tar -czf "$BACKUP_FILE" "$UPLOADS_DIR" 2>/dev/null

if [ -f "$BACKUP_FILE" ]; then
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}✅ Архив создан: $BACKUP_FILE (размер: $BACKUP_SIZE)${NC}"
else
    echo -e "${RED}❌ Ошибка при создании архива${NC}"
    exit 1
fi

# Удаление старых резервных копий
echo -e "${YELLOW}🗑️  Удаление старых резервных копий (оставляем последние $MAX_BACKUPS)...${NC}"

BACKUP_COUNT=$(ls -1 "$BACKUP_DIR"/uploads_*.tar.gz 2>/dev/null | wc -l | tr -d ' ')

if [ "$BACKUP_COUNT" -gt "$MAX_BACKUPS" ]; then
    ls -1t "$BACKUP_DIR"/uploads_*.tar.gz | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
    echo -e "${GREEN}✅ Удалено старых копий: $((BACKUP_COUNT - MAX_BACKUPS))${NC}"
fi

# Итоговая информация
echo ""
echo -e "${GREEN}✅ Резервное копирование завершено успешно!${NC}"
echo ""
echo -e "${YELLOW}📁 Созданный файл:${NC}"
echo "   - $BACKUP_FILE"
echo ""
echo -e "${YELLOW}📊 Всего резервных копий uploads:${NC}"
echo "   - Архивов: $(ls -1 "$BACKUP_DIR"/uploads_*.tar.gz 2>/dev/null | wc -l | tr -d ' ')"
echo ""
echo -e "${YELLOW}💡 Для восстановления используйте:${NC}"
echo "   tar -xzf $BACKUP_FILE"
echo ""
