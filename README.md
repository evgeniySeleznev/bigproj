# Week 1 - Микросервисы

Реализация трех микросервисов для системы управления заказами на космические корабли.

## Структура

```
bigproj/
├── inventory/    # Inventory Service (gRPC :50051)
├── payment/      # Payment Service (gRPC :50052)
├── order/        # Order Service (HTTP :8080)
└── shared/       # Контракты и сгенерированный код
```

## Запуск

```bash
# Терминал 1
cd inventory && go run cmd/server/main.go

# Терминал 2
cd payment && go run cmd/server/main.go

# Терминал 3
cd order && go run cmd/server/main.go
```

## API

**Order Service (HTTP :8080)**
- POST /api/v1/orders - создать заказ
- GET /api/v1/orders/{uuid} - получить заказ
- POST /api/v1/orders/{uuid}/pay - оплатить
- POST /api/v1/orders/{uuid}/cancel - отменить

**Inventory Service (gRPC :50051)**
- GetPart(uuid) - получить деталь
- ListParts(filter) - список деталей

**Payment Service (gRPC :50052)**
- PayOrder(req) - оплата заказа

## Тестирование

```bash
# Создать заказ
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_uuid":"user-1","part_uuids":["part-uuid-1"]}'

# Оплатить заказ
curl -X POST http://localhost:8080/api/v1/orders/{order_uuid}/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"CARD"}'
```

## Статус реализации

✅ Все три сервиса реализованы
✅ Интеграция через gRPC клиенты
✅ Соответствие контрактам
✅ Thread-safe хранилища
✅ Graceful shutdown





