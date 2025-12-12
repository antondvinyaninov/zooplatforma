#!/bin/bash

# Сравнение производительности с VK.com

echo "🔍 Сравнение производительности: ЗооПлатформа vs VK.com"
echo "========================================================"
echo ""

# Цвета
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Функция для тестирования
test_url() {
    local name=$1
    local url=$2
    local follow_redirects=$3
    
    if [ "$follow_redirects" = "yes" ]; then
        result=$(curl -L -o /dev/null -s -w "%{time_total},%{size_download}" "$url")
    else
        result=$(curl -o /dev/null -s -w "%{time_total},%{size_download}" "$url")
    fi
    
    IFS=',' read -r time size <<< "$result"
    
    # Конвертируем в миллисекунды
    time_ms=$(echo "$time * 1000" | bc -l)
    time_ms=$(printf "%.0f" $time_ms)
    
    # Конвертируем размер в KB
    size_kb=$(echo "scale=1; $size / 1024" | bc -l)
    
    printf "%-30s %6sms  %8sKB\n" "$name" "$time_ms" "$size_kb"
}

echo -e "${BLUE}📱 VK.com (Production, CDN, Кеширование):${NC}"
echo "--------------------------------------------"
test_url "Главная (vk.com)" "https://vk.com/" "yes"
test_url "Лента (vk.com/feed)" "https://vk.com/feed" "yes"

echo ""
echo -e "${GREEN}🐾 ЗооПлатформа (Local Dev):${NC}"
echo "--------------------------------------------"

if curl -s http://localhost:3000 > /dev/null 2>&1; then
    test_url "Главная (/)" "http://localhost:3000/" "no"
    test_url "Авторизация (/auth)" "http://localhost:3000/auth" "no"
    test_url "Организации (/org)" "http://localhost:3000/org" "no"
else
    echo "❌ Frontend не запущен"
fi

echo ""
echo "📊 Анализ:"
echo "--------------------------------------------"
echo "✅ VK.com - это production с многолетней оптимизацией:"
echo "   • CDN (Content Delivery Network)"
echo "   • Агрессивное кеширование"
echo "   • Минификация и сжатие"
echo "   • HTTP/2 и Server Push"
echo "   • Оптимизированные изображения"
echo "   • Code splitting"
echo ""
echo "🚀 ЗооПлатформа - это local dev режим:"
echo "   • Без кеширования"
echo "   • Без минификации"
echo "   • Без CDN"
echo "   • Hot reload включен"
echo "   • Source maps включены"
echo ""
echo "💡 Для честного сравнения нужно:"
echo "   1. npm run build (production build)"
echo "   2. Деплой на сервер с CDN"
echo "   3. Настроить кеширование"
echo "   4. Оптимизировать изображения"
echo ""
echo "🎯 Ожидаемое улучшение в production:"
echo "   • Скорость: 2-5x быстрее"
echo "   • Размер: 3-10x меньше"
echo "   • Кеширование: повторные загрузки <50ms"
