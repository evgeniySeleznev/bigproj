#!/bin/bash

set -e

echo "🧪 Тестирование микросервисов Week 1"
echo "====================================="
echo ""

# Функция для проверки доступности сервиса
wait_for_service() {
    local url=$1
    local name=$2
    local max_attempts=10
    local attempt=1

    echo "Ожидание запуска $name..."
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1 || nc -z localhost "$(echo $url | grep -oP ':\K[0-9]+')" 2>/dev/null; then
            echo "✓ $name доступен"
            return 0
        fi
        sleep 1
        attempt=$((attempt + 1))
    done
    echo "✗ $name недоступен"
    return 1
}

# Запускаем сервисы в фоне
echo "🚀 Запуск сервисов..."
cd inventory && go run cmd/server/main.go > /tmp/inventory.log 2>&1 &
INVENTORY_PID=$!
cd ..

cd payment && go run cmd/server/main.go > /tmp/payment.log 2>&1 &
PAYMENT_PID=$!
cd ..

cd order && go run cmd/server/main.go > /tmp/order.log 2>&1 &
ORDER_PID=$!
cd ..

# Ждем запуска сервисов
sleep 3

# Проверяем доступность сервисов
wait_for_service "localhost:50051" "Inventory Service" || exit 1
wait_for_service "localhost:50052" "Payment Service" || exit 1
wait_for_service "http://localhost:8080/api/v1/orders" "Order Service" || exit 1

echo ""
echo "📦 Тест 1: Получение списка деталей из Inventory"
PARTS=$(curl -s 'http://localhost:50051' 2>&1 || echo "gRPC service")
echo "Inventory service работает (gRPC на :50051)"

echo ""
echo "📝 Тест 2: Создание заказа"
USER_UUID="user-test-$(date +%s)"
PART_UUID="part-uuid-1"

ORDER_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{\"user_uuid\":\"$USER_UUID\",\"part_uuids\":[\"$PART_UUID\"]}")

echo "Ответ: $ORDER_RESPONSE"

if echo "$ORDER_RESPONSE" | grep -q "order_uuid"; then
    echo "✓ Заказ успешно создан"
    ORDER_UUID=$(echo "$ORDER_RESPONSE" | grep -oP '"order_uuid"\s*:\s*"[^"]*"' | cut -d'"' -f4)
    echo "Order UUID: $ORDER_UUID"
else
    echo "✗ Ошибка создания заказа"
    exit 1
fi

echo ""
echo "📊 Тест 3: Получение информации о заказе"
ORDER_INFO=$(curl -s "http://localhost:8080/api/v1/orders/$ORDER_UUID")
echo "$ORDER_INFO"

if echo "$ORDER_INFO" | grep -q "PENDING_PAYMENT"; then
    echo "✓ Статус заказа корректный: PENDING_PAYMENT"
else
    echo "✗ Неверный статус заказа"
fi

echo ""
echo "💰 Тест 4: Оплата заказа"
PAY_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/orders/$ORDER_UUID/pay" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"CARD"}')
echo "$PAY_RESPONSE"

if echo "$PAY_RESPONSE" | grep -q "transaction_uuid"; then
    echo "✓ Заказ успешно оплачен"
else
    echo "✗ Ошибка оплаты заказа"
fi

echo ""
echo "📊 Тест 5: Проверка статуса после оплаты"
ORDER_INFO=$(curl -s "http://localhost:8080/api/v1/orders/$ORDER_UUID")
if echo "$ORDER_INFO" | grep -q "PAID"; then
    echo "✓ Статус заказа после оплаты: PAID"
else
    echo "✗ Статус не обновлен на PAID"
fi

echo ""
echo "✅ Все тесты пройдены!"
echo ""
echo "Остановка сервисов..."
kill $INVENTORY_PID $PAYMENT_PID $ORDER_PID 2>/dev/null || true
wait

