#!/bin/bash

# Скрипт для тестирования производительности ЗооПлатформы

echo "🚀 Тестирование производительности ЗооПлатформы"
echo "================================================"
echo ""

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Функция для тестирования endpoint
test_endpoint() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    
    echo -n "📊 $name... "
    
    # Выполняем запрос и получаем метрики
    result=$(curl -o /dev/null -s -w "%{time_total},%{http_code},%{size_download}" -X $method "$url")
    
    IFS=',' read -r time status size <<< "$result"
    
    # Определяем цвет в зависимости от времени
    if (( $(echo "$time < 0.1" | bc -l) )); then
        color=$GREEN
        emoji="✅"
    elif (( $(echo "$time < 0.5" | bc -l) )); then
        color=$YELLOW
        emoji="⚠️"
    else
        color=$RED
        emoji="❌"
    fi
    
    # Форматируем время в миллисекунды
    time_ms=$(echo "$time * 1000" | bc -l)
    time_ms=$(printf "%.2f" $time_ms)
    
    # Форматируем размер
    if [ "$size" -gt 1024 ]; then
        size_kb=$(echo "scale=2; $size / 1024" | bc -l)
        size_str="${size_kb}KB"
    else
        size_str="${size}B"
    fi
    
    echo -e "${color}${emoji} ${time_ms}ms${NC} (${status}) [${size_str}]"
}

# Проверяем, запущен ли backend
if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend не запущен на порту 8080${NC}"
    echo "Запустите: cd backend && ./server"
    exit 1
fi

echo "Backend: ✅ Запущен"
echo ""

# Тестируем публичные endpoints
echo "📡 Публичные API endpoints:"
echo "----------------------------"
test_endpoint "Health Check      " "http://localhost:8080/api/health"
test_endpoint "Root              " "http://localhost:8080/"

echo ""

# Проверяем, запущен ли frontend
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "Frontend: ✅ Запущен"
    echo ""
    echo "🌐 Frontend страницы:"
    echo "----------------------------"
    test_endpoint "Главная (/)       " "http://localhost:3000/"
    test_endpoint "Авторизация       " "http://localhost:3000/auth"
    test_endpoint "Организации       " "http://localhost:3000/org"
else
    echo -e "${YELLOW}⚠️  Frontend не запущен на порту 3000${NC}"
    echo "Запустите: cd frontend && npm run dev"
fi

echo ""
echo "📈 Рекомендации:"
echo "----------------------------"
echo "✅ Отлично:  < 100ms"
echo "⚠️  Хорошо:   100-500ms"
echo "❌ Медленно: > 500ms"
echo ""
echo "💡 Для более детального анализа используйте:"
echo "   - Chrome DevTools (Network tab)"
echo "   - Lighthouse (Performance audit)"
echo "   - Next.js Analytics"
