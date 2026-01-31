#!/bin/bash

# Скрипт для проверки что все сервисы запущены

echo "🔍 Checking all services..."
echo ""

services=(
  "7000:Auth Service"
  "8000:Main Service"
  "8100:PetBase Service"
  "8200:Shelter Service"
  "8400:Owner Service"
  "8500:Volunteer Service"
  "8600:Clinic Service"
  "9000:Admin Service"
)

all_running=true

for service in "${services[@]}"; do
  port="${service%%:*}"
  name="${service##*:}"
  
  if curl -s "http://localhost:$port/api/health" > /dev/null 2>&1; then
    echo "✅ $name (port $port) - Running"
  else
    echo "❌ $name (port $port) - Not running"
    all_running=false
  fi
done

echo ""

if [ "$all_running" = true ]; then
  echo "🎉 All services are running!"
  exit 0
else
  echo "⚠️  Some services are not running. Check the logs."
  exit 1
fi
