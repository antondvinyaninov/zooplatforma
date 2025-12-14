#!/bin/bash

echo "🧹 Очистка кеша всех сервисов..."

# Main Frontend
if [ -d "main/frontend/.next" ]; then
  echo "  ✓ Очистка main/frontend/.next"
  rm -rf main/frontend/.next
fi

# Admin Frontend
if [ -d "admin/frontend/.next" ]; then
  echo "  ✓ Очистка admin/frontend/.next"
  rm -rf admin/frontend/.next
fi

# Mobile
if [ -d "mobile/.expo" ]; then
  echo "  ✓ Очистка mobile/.expo"
  rm -rf mobile/.expo
fi

if [ -d "mobile/node_modules/.cache" ]; then
  echo "  ✓ Очистка mobile/node_modules/.cache"
  rm -rf mobile/node_modules/.cache
fi

# Shared
if [ -d "shared/dist" ]; then
  echo "  ✓ Очистка shared/dist"
  rm -rf shared/dist
fi

echo ""
echo "✅ Кеш очищен!"
echo ""
echo "Теперь запустите сервисы:"
echo "  ./run"
