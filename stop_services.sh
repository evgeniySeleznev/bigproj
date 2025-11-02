#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="$BASE_DIR/services.pid"

if [ -f "$PID_FILE" ]; then
    echo "🛑 Останавливаем сервисы..."
    while read pid; do
        if kill -0 "$pid" 2>/dev/null; then
            echo "Останавливаем процесс $pid..."
            kill "$pid" 2>/dev/null
        fi
    done < "$PID_FILE"
    rm -f "$PID_FILE"
    echo "✅ Все сервисы остановлены"
else
    echo "⚠️  Сервисы не запущены (файл $PID_FILE не найден)"
fi
