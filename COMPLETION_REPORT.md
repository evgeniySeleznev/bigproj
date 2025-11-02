# Week 1 - Отчет о выполнении

## ✅ Выполненные требования

### 1. OrderService (HTTP API)
**Требование:** HTTP API строго следуя OpenAPI-контракту

**Реализовано:**
- ✅ POST /api/v1/orders - создание заказа
- ✅ GET /api/v1/orders/{order_uuid} - получение заказа
- ✅ POST /api/v1/orders/{order_uuid}/pay - оплата заказа
- ✅ POST /api/v1/orders/{order_uuid}/cancel - отмена заказа

**Детали реализации:**
- Интеграция с InventoryService для получения деталей через ListParts
- Проверка существования всех деталей
- Подсчет total_price из деталей
- Генерация order_uuid
- Сохранение заказа со статусом PENDING_PAYMENT
- Интеграция с PaymentService для оплаты
- Обновление статуса на PAID после оплаты
- Проверка статуса при отмене (409 Conflict если PAID)
- Хранилище в памяти с sync.RWMutex

### 2. InventoryService (gRPC API)
**Требование:** gRPC API по контракту inventory_service_contracts.md

**Реализовано:**
- ✅ GetPart(context, *GetPartRequest) (*GetPartResponse)
- ✅ ListParts(context, *ListPartsRequest) (*ListPartsResponse)

**Детали реализации:**
- Хранилище map[string]*Part с sync.RWMutex
- Инициализация с 4 тестовыми деталями:
  - Ion Engine Model X1 (ENGINE)
  - Liquid Hydrogen Tank 500L (FUEL)
  - Observation Window 50cm (PORTHOLE)
  - Solar Wing Panel 4m (WING)
- GetPart: поиск по UUID, возврат ошибки NotFound если не найден
- ListParts: фильтрация по PartsFilter
  - Логическое ИЛИ внутри поля
  - Логическое И между полями
  - Фильтрация по: uuids, names, categories, manufacturer_countries, tags
- gRPC сервер на порту 50051
- Graceful shutdown

### 3. PaymentService (gRPC API)
**Требование:** gRPC API по контракту payment_service_contracts.md

**Реализовано:**
- ✅ PayOrder(context, *PayOrderRequest) (*PayOrderResponse)

**Детали реализации:**
- Генерация transaction_uuid через uuid.New()
- Логирование в консоль: "Оплата прошла успешно, transaction_uuid: ..."
- Возврат transaction_uuid
- Stateless сервис
- gRPC сервер на порту 50052
- Graceful shutdown

### 4. Интеграция через gRPC клиенты
**Требование:** В OrderService через gRPC-клиенты интегрировать InventoryService и PaymentService

**Реализовано:**
- ✅ OrderService создает gRPC клиенты для обоих сервисов
- ✅ При создании заказа: вызов InventoryService.ListParts
- ✅ При оплате заказа: вызов PaymentService.PayOrder
- ✅ Обработка ошибок от gRPC сервисов (Bad Gateway)
- ✅ Подключение к localhost:50051 и localhost:50052

## 📁 Созданные файлы

### Структура проекта
```
myWeek_1/
├── go.work ✅
├── buf.work.yaml ✅
├── Taskfile.yml ✅
├── package.json ✅
├── .golangci.yml ✅
├── inventory/
│   ├── cmd/server/main.go ✅
│   └── go.mod ✅
├── payment/
│   ├── cmd/server/main.go ✅
│   └── go.mod ✅
├── order/
│   ├── cmd/server/main.go ✅
│   └── go.mod ✅
└── shared/
    ├── go.mod ✅
    ├── proto/
    │   ├── inventory/v1/inventory.proto ✅
    │   └── payment/v1/payment.proto ✅
    ├── api/
    │   └── order/v1/ (все компоненты) ✅
    └── pkg/proto/
        ├── inventory/v1/ (сгенерированный код) ✅
        └── payment/v1/ (сгенерированный код) ✅
```

## 🧪 Тестирование

### Проверка компиляции
```bash
✓ inventory compiled successfully
✓ payment compiled successfully
✓ order compiled successfully
```

### Для запуска и тестирования:

**Терминал 1 - Inventory Service:**
```bash
cd /home/evgeniyalter/micro/week1/myWeek_1/inventory
go run cmd/server/main.go
```

**Терминал 2 - Payment Service:**
```bash
cd /home/evgeniyalter/micro/week1/myWeek_1/payment
go run cmd/server/main.go
```

**Терминал 3 - Order Service:**
```bash
cd /home/evgeniyalter/micro/week1/myWeek_1/order
go run cmd/server/main.go
```

**Терминал 4 - Тестирование:**
```bash
# Создать заказ
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_uuid":"user-1","part_uuids":["part-uuid-1"]}'

# Получить заказ (замените {order_uuid})
curl http://localhost:8080/api/v1/orders/{order_uuid}

# Оплатить заказ
curl -X POST http://localhost:8080/api/v1/orders/{order_uuid}/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"CARD"}'

# Отменить заказ (только если не оплачен)
curl -X POST http://localhost:8080/api/v1/orders/{order_uuid}/cancel
```

## ✅ Соответствие требованиям hw.md

1. ✅ Реализован HTTP API для OrderService согласно контракту
2. ✅ Реализован gRPC API для InventoryService согласно контракту
3. ✅ Реализован gRPC API для PaymentService согласно контракту
4. ✅ OrderService интегрирован с Inventory и Payment через gRPC клиенты
5. ✅ Использован монорепозиторий с Go Workspaces (go.work)
6. ✅ Проект структурирован согласно рекомендациям
7. ✅ Все контракты находятся в shared/
8. ✅ Логика реализована в main.go (будет вынесена в слои позже)

## 🎯 Готово к использованию

Все сервисы готовы к запуску и тестированию. Следующий шаг - запуск всех трех сервисов и проверка работы через HTTP запросы.








