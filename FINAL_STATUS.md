# Week 1 - Итоговый статус

## ✅ Выполнено

### 1. Инфраструктура
- go.work, buf.work.yaml, Taskfile.yml, .golangci.yml созданы
- Проект в `myWeek_1/`

### 2. Proto контракты
- `shared/proto/inventory/v1/inventory.proto` ✅
- `shared/proto/payment/v1/payment.proto` ✅
- Сгенерирован код в `shared/pkg/proto/`

### 3. Сервисы реализованы
- Inventory Service: gRPC на :50051
  - GetPart(), ListParts()
  - 4 тестовые детали
- Payment Service: gRPC на :50052  
  - PayOrder()
- Order Service: HTTP на :8080
  - POST /api/v1/orders
  - GET /api/v1/orders/{uuid}
  - POST /api/v1/orders/{uuid}/pay
  - POST /api/v1/orders/{uuid}/cancel

## 📁 Файлы
- `inventory/cmd/server/main.go` ✅
- `payment/cmd/server/main.go` ✅
- `order/cmd/server/main.go` ✅

## 🚀 Запуск
```bash
# Терминал 1
cd myWeek_1/inventory && go run cmd/server/main.go

# Терминал 2
cd myWeek_1/payment && go run cmd/server/main.go

# Терминал 3
cd myWeek_1/order && go run cmd/server/main.go
```

## 🧪 Тест
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_uuid":"user-1","part_uuids":["part-uuid-1"]}'
```

