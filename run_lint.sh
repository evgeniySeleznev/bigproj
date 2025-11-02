#!/bin/bash
# Скрипт для запуска линтера напрямую

set -e

cd "$(dirname "$0")"

echo "🔍 Запуск golangci-lint для всех модулей..."
echo ""

ERRORS=0

echo "📦 Проверка inventory..."
if ./bin/golangci-lint run ./inventory/... --config=.golangci.yml; then
    echo "✅ inventory: OK"
else
    echo "❌ inventory: Ошибки найдены"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "📦 Проверка payment..."
if ./bin/golangci-lint run ./payment/... --config=.golangci.yml; then
    echo "✅ payment: OK"
else
    echo "❌ payment: Ошибки найдены"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "📦 Проверка order..."
if ./bin/golangci-lint run ./order/... --config=.golangci.yml; then
    echo "✅ order: OK"
else
    echo "❌ order: Ошибки найдены"
    ERRORS=$((ERRORS + 1))
fi

echo ""
if [ $ERRORS -eq 0 ]; then
    echo "✅ Все модули прошли проверку линтера!"
    exit 0
else
    echo "❌ Найдены ошибки в $ERRORS модуле(ях)"
    exit 1
fi

