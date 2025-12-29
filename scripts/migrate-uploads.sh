#!/bin/bash

# ============================================
# Скрипт миграции файлов uploads
# Переносит файлы из main/backend/uploads/ в корневую uploads/
# ============================================

set -e  # Остановка при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔄 Миграция файлов uploads${NC}"
echo ""

OLD_UPLOADS="main/backend/uploads"
NEW_UPLOADS="uploads"

# Проверка существования старой папки
if [ ! -d "$OLD_UPLOADS" ]; then
    echo -e "${YELLOW}⚠️  Папка $OLD_UPLOADS не найдена. Миграция не требуется.${NC}"
    exit 0
fi

# Проверка существования новой папки
if [ ! -d "$NEW_UPLOADS" ]; then
    echo -e "${YELLOW}📁 Создание папки $NEW_UPLOADS${NC}"
    mkdir -p "$NEW_UPLOADS"
fi

# Подсчёт файлов для миграции
FILE_COUNT=$(find "$OLD_UPLOADS" -type f ! -name ".gitkeep" | wc -l | tr -d ' ')

if [ "$FILE_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✅ Нет файлов для миграции${NC}"
    exit 0
fi

echo -e "${YELLOW}📊 Найдено файлов для миграции: $FILE_COUNT${NC}"
echo ""

# Показываем структуру
echo -e "${YELLOW}📁 Структура файлов:${NC}"
tree "$OLD_UPLOADS" 2>/dev/null || find "$OLD_UPLOADS" -type f ! -name ".gitkeep" | head -10
echo ""

read -p "Продолжить миграцию? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}❌ Миграция отменена${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}🔄 Копирование файлов...${NC}"

# Копируем все файлы, сохраняя структуру
rsync -av --progress "$OLD_UPLOADS/" "$NEW_UPLOADS/" --exclude=".gitkeep"

# Проверяем результат
NEW_FILE_COUNT=$(find "$NEW_UPLOADS" -type f ! -name ".gitkeep" | wc -l | tr -d ' ')

echo ""
if [ "$NEW_FILE_COUNT" -eq "$FILE_COUNT" ]; then
    echo -e "${GREEN}✅ Все файлы успешно скопированы ($NEW_FILE_COUNT файлов)${NC}"
    echo ""
    
    # Предлагаем удалить старую папку
    read -p "Удалить старую папку $OLD_UPLOADS? (yes/no): " DELETE_OLD
    
    if [ "$DELETE_OLD" = "yes" ]; then
        echo -e "${YELLOW}🗑️  Удаление старой папки...${NC}"
        rm -rf "$OLD_UPLOADS"
        echo -e "${GREEN}✅ Старая папка удалена${NC}"
    else
        echo -e "${YELLOW}⚠️  Старая папка сохранена. Удалите вручную после проверки:${NC}"
        echo "   rm -rf $OLD_UPLOADS"
    fi
else
    echo -e "${RED}❌ Ошибка: количество файлов не совпадает${NC}"
    echo "   Ожидалось: $FILE_COUNT"
    echo "   Скопировано: $NEW_FILE_COUNT"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Миграция завершена!${NC}"
echo ""
echo -e "${YELLOW}📁 Новая структура:${NC}"
tree "$NEW_UPLOADS" 2>/dev/null || ls -la "$NEW_UPLOADS"
echo ""
echo -e "${YELLOW}💡 Следующие шаги:${NC}"
echo "   1. Перезапустите backend: cd main/backend && go run main.go"
echo "   2. Проверьте загрузку файлов"
echo "   3. Проверьте отображение аватаров и обложек"
echo ""
