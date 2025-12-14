#!/bin/bash

echo "🐾 Starting ЗооБаза..."

# Clean cache
echo "🧹 Cleaning cache..."
rm -rf frontend/.next 2>/dev/null || true
echo "✅ Cache cleared"

# Kill processes on ports
echo "🧹 Cleaning up ports..."
lsof -ti:4100 | xargs kill -9 2>/dev/null || true
lsof -ti:8100 | xargs kill -9 2>/dev/null || true
echo "✅ Ports cleared"

# Start backend
echo "🔧 Starting ЗооБаза Backend on port 8100..."
(cd backend && go run main.go) &
BACKEND_PID=$!

# Wait for backend
sleep 2

# Start frontend
echo "⚛️  Starting ЗооБаза Frontend on port 4100..."
(cd frontend && npm run dev) &
FRONTEND_PID=$!

echo ""
echo "✅ ЗооБаза запущена!"
echo ""
echo "   🌐 Frontend:  http://localhost:4100"
echo "   🔧 Backend:   http://localhost:8100"
echo ""
echo "Press Ctrl+C to stop"

# Wait for Ctrl+C
trap "echo ''; echo '🛑 Stopping ЗооБаза...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT
wait
