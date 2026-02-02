#!/bin/bash

# Скрипт для логирования HTTP запросов с деталями
# Добавляет логирование в каждый backend сервис

echo "📊 Настройка логирования HTTP запросов..."
echo ""

# Создаем временную директорию для логов
LOG_DIR="logs/requests"
mkdir -p "$LOG_DIR"

echo "Логи будут сохраняться в: $LOG_DIR"
echo ""
echo "Сервисы:"
echo "  Main Backend (8000)     -> $LOG_DIR/main.log"
echo "  Admin Backend (9000)    -> $LOG_DIR/admin.log"
echo "  PetBase Backend (8100)  -> $LOG_DIR/petbase.log"
echo "  Shelter Backend (8200)  -> $LOG_DIR/shelter.log"
echo "  Owner Backend (8400)    -> $LOG_DIR/owner.log"
echo "  Volunteer Backend (8500)-> $LOG_DIR/volunteer.log"
echo "  Clinic Backend (8600)   -> $LOG_DIR/clinic.log"
echo "  Auth Service (7100)     -> $LOG_DIR/auth.log"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Функция для мониторинга логов
monitor_logs() {
    # Очищаем старые логи
    > "$LOG_DIR/main.log"
    > "$LOG_DIR/admin.log"
    > "$LOG_DIR/petbase.log"
    > "$LOG_DIR/shelter.log"
    > "$LOG_DIR/owner.log"
    > "$LOG_DIR/volunteer.log"
    > "$LOG_DIR/clinic.log"
    > "$LOG_DIR/auth.log"
    
    echo "🔍 Мониторинг запросов (нажмите Ctrl+C для остановки)..."
    echo ""
    
    # Мониторим все логи одновременно с цветами
    tail -f \
        "$LOG_DIR/main.log" \
        "$LOG_DIR/admin.log" \
        "$LOG_DIR/petbase.log" \
        "$LOG_DIR/shelter.log" \
        "$LOG_DIR/owner.log" \
        "$LOG_DIR/volunteer.log" \
        "$LOG_DIR/clinic.log" \
        "$LOG_DIR/auth.log" 2>/dev/null | \
    while read line; do
        timestamp=$(date '+%H:%M:%S')
        
        # Определяем цвет по методу
        if echo "$line" | grep -q "GET"; then
            echo -e "\033[0;32m[$timestamp] $line\033[0m"  # Зеленый
        elif echo "$line" | grep -q "POST"; then
            echo -e "\033[0;34m[$timestamp] $line\033[0m"  # Синий
        elif echo "$line" | grep -q "PUT"; then
            echo -e "\033[0;33m[$timestamp] $line\033[0m"  # Желтый
        elif echo "$line" | grep -q "DELETE"; then
            echo -e "\033[0;31m[$timestamp] $line\033[0m"  # Красный
        else
            echo "[$timestamp] $line"
        fi
    done
}

# Запускаем мониторинг
monitor_logs
