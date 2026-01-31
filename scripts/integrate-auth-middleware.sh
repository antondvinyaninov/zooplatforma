#!/bin/bash

echo "🔐 Integrating centralized Auth middleware for all services..."

# Список сервисов для обновления
SERVICES=("admin" "clinic" "shelter" "owner" "volunteer" "petbase")

for service in "${SERVICES[@]}"; do
    echo ""
    echo "📦 Processing $service..."
    
    # Переходим в директорию backend
    cd "$service/backend" || continue
    
    # Добавляем зависимость на pkg/middleware
    echo "   Adding pkg/middleware dependency..."
    if ! grep -q "github.com/zooplatforma/pkg/middleware" go.mod; then
        # Добавляем replace directive
        if ! grep -q "replace github.com/zooplatforma/pkg/middleware" go.mod; then
            echo "" >> go.mod
            echo "replace github.com/zooplatforma/pkg/middleware => ../../pkg/middleware" >> go.mod
        fi
        
        # Добавляем require
        if ! grep -q "github.com/zooplatforma/pkg/middleware" go.mod; then
            # Находим последнюю строку require и добавляем после нее
            sed -i '' '/^require (/a\
	github.com/zooplatforma/pkg/middleware v0.0.0
' go.mod
        fi
    fi
    
    # Запускаем go mod tidy
    echo "   Running go mod tidy..."
    go mod tidy
    
    echo "   ✅ $service updated"
    
    # Возвращаемся в корень
    cd ../..
done

echo ""
echo "✅ All services updated!"
echo ""
echo "📝 Next steps:"
echo "1. Update main.go in each service to import pkg/middleware"
echo "2. Replace local AuthMiddleware with middleware.AuthMiddleware"
echo "3. Update handlers to use middleware.GetUserID(r)"
echo "4. Test each service"
echo ""
echo "See AUTH_SERVICE_READY.md for detailed instructions"
