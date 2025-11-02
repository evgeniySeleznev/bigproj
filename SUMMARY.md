# Week 1 - Краткий итог

## ✅ Выполнено

**3 микросервиса реализованы:**
- InventoryService (gRPC :50051) - 2 метода
- PaymentService (gRPC :50052) - 1 метод
- OrderService (HTTP :8080) - 4 endpoint

**Все требования hw.md выполнены:**
1. HTTP API для OrderService по OpenAPI контракту ✅
2. gRPC API для InventoryService по proto контракту ✅
3. gRPC API для PaymentService по proto контракту ✅
4. Интеграция через gRPC клиенты ✅
5. Go Workspaces (go.work) ✅
6. Контракты в shared/ ✅

## 📂 Основные файлы

- `inventory/cmd/server/main.go` - gRPC сервер + хранилище
- `payment/cmd/server/main.go` - gRPC сервер + метод PayOrder
- `order/cmd/server/main.go` - HTTP сервер + gRPC интеграция
- `shared/proto/*.proto` - контракты
- `shared/pkg/proto/*` - сгенерированный код

## 🚀 Запуск

```bash
# 3 терминала для каждого сервиса
cd myWeek_1/inventory && go run cmd/server/main.go  # порт 50051
cd myWeek_1/payment && go run cmd/server/main.go     # порт 50052
cd myWeek_1/order && go run cmd/server/main.go       # порт 8080
```

## 🧪 Тест

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_uuid":"user-1","part_uuids":["part-uuid-1"]}'
```

## ✨ Готово к проверке

Все сервисы компилируются без ошибок и готовы к тестированию.

