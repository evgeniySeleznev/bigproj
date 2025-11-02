#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="$BASE_DIR/services.pid"

# Функция для остановки сервисов
stop_services() {
    if [ -f "$PID_FILE" ]; then
        echo -e "${YELLOW}Останавливаем сервисы...${NC}"
        while read pid; do
            if kill -0 "$pid" 2>/dev/null; then
                kill "$pid" 2>/dev/null
            fi
        done < "$PID_FILE"
        rm -f "$PID_FILE"
        echo -e "${GREEN}Сервисы остановлены${NC}"
    else
        echo -e "${YELLOW}Сервисы не запущены${NC}"
    fi
    exit 0
}

# Обработка сигналов для корректного завершения
trap stop_services SIGINT SIGTERM

# Остановка существующих сервисов
stop_services > /dev/null 2>&1

echo -e "${BLUE}🚀 Запускаем все сервисы...${NC}"

# Запуск Inventory Service (gRPC :50051)
cd "$BASE_DIR/inventory" || exit 1
go run cmd/server/main.go > /tmp/inventory.log 2>&1 &
INVENTORY_PID=$!
echo "$INVENTORY_PID" >> "$PID_FILE"
echo -e "${GREEN}✅ Inventory Service запущен (PID: $INVENTORY_PID) на порту :50051${NC}"
echo "   Логи: tail -f /tmp/inventory.log"

# Небольшая задержка для запуска первого сервиса
sleep 1

# Запуск Payment Service (gRPC :50052)
cd "$BASE_DIR/payment" || exit 1
go run cmd/server/main.go > /tmp/payment.log 2>&1 &
PAYMENT_PID=$!
echo "$PAYMENT_PID" >> "$PID_FILE"
echo -e "${GREEN}✅ Payment Service запущен (PID: $PAYMENT_PID) на порту :50052${NC}"
echo "   Логи: tail -f /tmp/payment.log"

sleep 1

# Запуск Order Service (HTTP :8080)
cd "$BASE_DIR/order" || exit 1
go run cmd/server/main.go > /tmp/order.log 2>&1 &
ORDER_PID=$!
echo "$ORDER_PID" >> "$PID_FILE"
echo -e "${GREEN}✅ Order Service запущен (PID: $ORDER_PID) на порту :8080${NC}"
echo "   Логи: tail -f /tmp/order.log"

sleep 2

echo ""
echo -e "${BLUE}📋 Все сервисы запущены!${NC}"
echo -e "   ${GREEN}Inventory Service${NC}: gRPC :50051"
echo -e "   ${GREEN}Payment Service${NC}:   gRPC :50052"
echo -e "   ${GREEN}Order Service${NC}:     HTTP :8080"
echo ""
echo -e "${YELLOW}Для остановки нажмите Ctrl+C или выполните: ./stop_services.sh${NC}"
echo ""

# Ожидание завершения
wait
