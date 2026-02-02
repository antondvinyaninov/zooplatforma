#!/bin/bash

# Простой скрипт для просмотра активных соединений к сервисам
# Показывает количество запросов к каждому сервису

echo "📊 Мониторинг активных соединений к сервисам"
echo "Обновление каждые 2 секунды. Нажмите Ctrl+C для остановки"
echo ""

while true; do
    clear
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  📊 Активные соединения к сервисам"
    echo "  Время: $(date '+%H:%M:%S')"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Функция для подсчета соединений
    count_connections() {
        port=$1
        name=$2
        count=$(netstat -an 2>/dev/null | grep ".$port " | grep ESTABLISHED | wc -l | tr -d ' ')
        
        if [ "$count" -gt 0 ]; then
            echo "  🟢 $name (port $port): $count активных соединений"
        else
            echo "  ⚪ $name (port $port): нет соединений"
        fi
    }
    
    echo "Frontend сервисы:"
    count_connections 3000 "Main Frontend    "
    count_connections 4000 "Admin Frontend   "
    count_connections 4100 "PetBase Frontend "
    count_connections 5100 "Shelter Frontend "
    count_connections 6100 "Owner Frontend   "
    count_connections 6200 "Volunteer Frontend"
    count_connections 6300 "Clinic Frontend  "
    
    echo ""
    echo "Backend сервисы:"
    count_connections 8000 "Main Backend     "
    count_connections 9000 "Admin Backend    "
    count_connections 8100 "PetBase Backend  "
    count_connections 8200 "Shelter Backend  "
    count_connections 8400 "Owner Backend    "
    count_connections 8500 "Volunteer Backend"
    count_connections 8600 "Clinic Backend   "
    count_connections 7100 "Auth Service     "
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Общая статистика:"
    
    # Подсчет общего количества соединений
    total_frontend=$(netstat -an 2>/dev/null | grep -E "\.(3000|4000|4100|5100|6100|6200|6300) " | grep ESTABLISHED | wc -l | tr -d ' ')
    total_backend=$(netstat -an 2>/dev/null | grep -E "\.(8000|9000|8100|8200|8400|8500|8600|7100) " | grep ESTABLISHED | wc -l | tr -d ' ')
    total=$((total_frontend + total_backend))
    
    echo "  Frontend: $total_frontend соединений"
    echo "  Backend:  $total_backend соединений"
    echo "  Всего:    $total соединений"
    
    echo ""
    echo "Нажмите Ctrl+C для остановки..."
    
    sleep 2
done
