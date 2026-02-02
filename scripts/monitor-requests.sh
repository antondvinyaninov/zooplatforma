#!/bin/bash

# Скрипт для мониторинга HTTP запросов ко всем сервисам
# Показывает какие запросы идут к каким портам

echo "🔍 Мониторинг HTTP запросов к сервисам..."
echo "Нажмите Ctrl+C для остановки"
echo ""
echo "Порты сервисов:"
echo "  3000 - Main Frontend"
echo "  8000 - Main Backend"
echo "  4000 - Admin Frontend"
echo "  9000 - Admin Backend"
echo "  4100 - PetBase Frontend"
echo "  8100 - PetBase Backend"
echo "  5100 - Shelter Frontend"
echo "  8200 - Shelter Backend"
echo "  6100 - Owner Frontend"
echo "  8400 - Owner Backend"
echo "  6200 - Volunteer Frontend"
echo "  8500 - Volunteer Backend"
echo "  6300 - Clinic Frontend"
echo "  8600 - Clinic Backend"
echo "  7100 - Auth Service"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Функция для получения имени сервиса по порту
get_service_name() {
    case $1 in
        3000) echo "Main Frontend" ;;
        8000) echo "Main Backend" ;;
        4000) echo "Admin Frontend" ;;
        9000) echo "Admin Backend" ;;
        4100) echo "PetBase Frontend" ;;
        8100) echo "PetBase Backend" ;;
        5100) echo "Shelter Frontend" ;;
        8200) echo "Shelter Backend" ;;
        6100) echo "Owner Frontend" ;;
        8400) echo "Owner Backend" ;;
        6200) echo "Volunteer Frontend" ;;
        8500) echo "Volunteer Backend" ;;
        6300) echo "Clinic Frontend" ;;
        8600) echo "Clinic Backend" ;;
        7100) echo "Auth Service" ;;
        *) echo "Unknown:$1" ;;
    esac
}

# Цвета для разных типов запросов
COLOR_GET="\033[0;32m"      # Зеленый
COLOR_POST="\033[0;34m"     # Синий
COLOR_PUT="\033[0;33m"      # Желтый
COLOR_DELETE="\033[0;31m"   # Красный
COLOR_RESET="\033[0m"

# Мониторим сетевые соединения
if command -v lsof &> /dev/null; then
    # macOS/Linux с lsof
    while true; do
        # Получаем активные соединения на наших портах
        lsof -iTCP:3000,8000,4000,9000,4100,8100,5100,8200,6100,8400,6200,8500,6300,8600,7100 -sTCP:ESTABLISHED -n -P 2>/dev/null | \
        grep -v "COMMAND" | \
        while read line; do
            # Парсим вывод lsof
            port=$(echo "$line" | awk '{print $9}' | grep -o ':[0-9]*' | cut -d: -f2)
            if [ ! -z "$port" ]; then
                service=$(get_service_name $port)
                timestamp=$(date '+%H:%M:%S')
                echo "[$timestamp] 🔗 Connection to $service (port $port)"
            fi
        done
        sleep 1
    done
else
    echo "❌ lsof не найден. Устанавливаем альтернативный метод..."
    echo ""
    echo "Для macOS используйте: brew install lsof"
    echo ""
    echo "Альтернатива: мониторинг через netstat..."
    
    while true; do
        netstat -an | grep ESTABLISHED | grep -E ":(3000|8000|4000|9000|4100|8100|5100|8200|6100|8400|6200|8500|6300|8600|7100)" | \
        while read line; do
            port=$(echo "$line" | grep -o '\.[0-9]*' | head -1 | cut -d. -f2)
            if [ ! -z "$port" ]; then
                service=$(get_service_name $port)
                timestamp=$(date '+%H:%M:%S')
                echo "[$timestamp] 🔗 Connection to $service (port $port)"
            fi
        done
        sleep 1
    done
fi
