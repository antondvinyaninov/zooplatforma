#!/bin/bash

# Скрипт для измерения Web Vitals с помощью Lighthouse CLI

echo "🎯 Измерение Web Vitals с помощью Lighthouse"
echo "=============================================="
echo ""

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Проверяем, установлен ли Lighthouse
if ! command -v lighthouse &> /dev/null; then
    echo -e "${RED}❌ Lighthouse не установлен${NC}"
    echo ""
    echo "Установите Lighthouse CLI:"
    echo -e "${BLUE}npm install -g lighthouse${NC}"
    echo ""
    echo "Или используйте Chrome DevTools:"
    echo "1. Откройте DevTools (F12)"
    echo "2. Вкладка 'Lighthouse'"
    echo "3. Выберите 'Performance'"
    echo "4. Нажмите 'Analyze page load'"
    exit 1
fi

# Проверяем, запущен ли frontend
if ! curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${RED}❌ Frontend не запущен на порту 3000${NC}"
    echo "Запустите: cd frontend && npm run dev"
    exit 1
fi

echo -e "${GREEN}✅ Frontend запущен${NC}"
echo ""

# Создаем папку для отчетов
mkdir -p tests/reports

# Функция для запуска Lighthouse
run_lighthouse() {
    local url=$1
    local name=$2
    
    echo -e "${BLUE}📊 Тестирование: $name${NC}"
    echo "URL: $url"
    echo ""
    
    # Запускаем Lighthouse
    lighthouse "$url" \
        --only-categories=performance \
        --output=json \
        --output-path="tests/reports/${name}.json" \
        --chrome-flags="--headless" \
        --quiet
    
    # Читаем результаты
    if [ -f "tests/reports/${name}.json" ]; then
        # Извлекаем метрики с помощью jq (если установлен)
        if command -v jq &> /dev/null; then
            score=$(jq -r '.categories.performance.score * 100' "tests/reports/${name}.json")
            fcp=$(jq -r '.audits["first-contentful-paint"].displayValue' "tests/reports/${name}.json")
            lcp=$(jq -r '.audits["largest-contentful-paint"].displayValue' "tests/reports/${name}.json")
            tbt=$(jq -r '.audits["total-blocking-time"].displayValue' "tests/reports/${name}.json")
            cls=$(jq -r '.audits["cumulative-layout-shift"].displayValue' "tests/reports/${name}.json")
            
            # Определяем цвет для score
            if (( $(echo "$score >= 90" | bc -l) )); then
                color=$GREEN
                emoji="✅"
            elif (( $(echo "$score >= 50" | bc -l) )); then
                color=$YELLOW
                emoji="⚠️"
            else
                color=$RED
                emoji="❌"
            fi
            
            echo -e "${color}${emoji} Performance Score: ${score}/100${NC}"
            echo "   FCP: $fcp"
            echo "   LCP: $lcp"
            echo "   TBT: $tbt"
            echo "   CLS: $cls"
            echo ""
        else
            echo "✅ Отчет сохранен: tests/reports/${name}.json"
            echo "   Установите jq для отображения метрик: brew install jq"
            echo ""
        fi
    fi
}

# Тестируем страницы
echo "🚀 Запуск тестов..."
echo ""

run_lighthouse "http://localhost:3000/" "homepage"
run_lighthouse "http://localhost:3000/auth" "auth"
run_lighthouse "http://localhost:3000/org" "organizations"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Результаты сохранены в tests/reports/"
echo ""
echo "📖 Расшифровка метрик:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Performance Score:"
echo "  90-100: ✅ Отлично"
echo "  50-89:  ⚠️  Требует улучшения"
echo "  0-49:   ❌ Плохо"
echo ""
echo "FCP (First Contentful Paint):"
echo "  <1.8s:  ✅ Оптимально"
echo "  1.8-3s: ⚠️  Средне"
echo "  >3s:    ❌ Медленно"
echo ""
echo "LCP (Largest Contentful Paint):"
echo "  <2.5s:  ✅ Хорошо"
echo "  2.5-4s: ⚠️  Требует улучшения"
echo "  >4s:    ❌ Плохо"
echo ""
echo "TBT (Total Blocking Time):"
echo "  <200ms: ✅ Хорошо"
echo "  200-600ms: ⚠️  Средне"
echo "  >600ms: ❌ Плохо"
echo ""
echo "CLS (Cumulative Layout Shift):"
echo "  <0.1:   ✅ Хорошо"
echo "  0.1-0.25: ⚠️  Средне"
echo "  >0.25:  ❌ Плохо"
echo ""
echo "💡 Для детального анализа откройте:"
echo "   tests/reports/*.json"
echo ""
echo "🌐 Или используйте Chrome DevTools → Lighthouse"
