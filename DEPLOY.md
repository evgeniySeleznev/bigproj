# Инструкция по загрузке в GitHub

## Создание репозитория на GitHub

1. Перейдите на https://github.com/new
2. Название репозитория: `bigproj`
3. Владелец: `evgeniyseleznev`
4. Описание: "Микросервисы для системы управления заказами на космические корабли"
5. Выберите Public или Private
6. **НЕ** добавляйте README, .gitignore или лицензию (мы добавим их вручную)
7. Нажмите "Create repository"

## Команды для загрузки

```bash
cd "/home/evgeniyalter/micro/week1/ver 0.0.1"

# Инициализация git (если еще не инициализирован)
git init

# Добавление всех файлов
git add .

# Первый коммит
git commit -m "Initial commit: Week 1 - Микросервисы Order, Inventory, Payment"

# Добавление remote репозитория (замените YOUR_TOKEN на ваш GitHub токен, если используете HTTPS)
git remote add origin https://github.com/evgeniyseleznev/bigproj.git

# Или через SSH (если настроен):
# git remote add origin git@github.com:evgeniyseleznev/bigproj.git

# Проверка remote
git remote -v

# Загрузка в репозиторий
git branch -M main
git push -u origin main
```

## Если нужно обновить токен доступа

Если используете HTTPS и нужна аутентификация:
1. Создайте Personal Access Token в GitHub Settings → Developer settings → Personal access tokens
2. Используйте токен вместо пароля при push

## Структура репозитория

После загрузки структура будет:
```
bigproj/
├── inventory/        # Inventory Service (gRPC)
├── order/            # Order Service (HTTP)
├── payment/          # Payment Service (gRPC)
├── shared/           # Общие контракты и протобуфы
├── go.work           # Go Workspace
├── README.md
└── ...
```

