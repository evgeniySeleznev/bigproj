#!/bin/bash

# Скрипт для загрузки проекта в GitHub репозиторий
# Использование: ./deploy.sh

set -e

REPO_URL="https://github.com/evgeniyseleznev/bigproj.git"
CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$CURRENT_DIR"

echo "📦 Подготовка к загрузке в GitHub..."

# Проверка наличия git
if ! command -v git &> /dev/null; then
    echo "❌ Git не установлен. Установите git и повторите попытку."
    exit 1
fi

# Инициализация git (если еще не инициализирован)
if [ ! -d ".git" ]; then
    echo "🔧 Инициализация git репозитория..."
    git init
fi

# Проверка remote
if git remote get-url origin &> /dev/null; then
    CURRENT_REMOTE=$(git remote get-url origin)
    echo "📡 Текущий remote: $CURRENT_REMOTE"
    read -p "Изменить remote на $REPO_URL? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git remote set-url origin "$REPO_URL" || git remote add origin "$REPO_URL"
    fi
else
    echo "📡 Добавление remote репозитория..."
    git remote add origin "$REPO_URL"
fi

# Добавление всех файлов
echo "➕ Добавление файлов в git..."
git add .

# Проверка изменений
if git diff --cached --quiet && git diff --quiet; then
    echo "ℹ️  Нет изменений для коммита"
else
    # Создание коммита
    echo "💾 Создание коммита..."
    git commit -m "Initial commit: Week 1 - Микросервисы Order, Inventory, Payment

- Реализован Order Service (HTTP API)
- Реализован Inventory Service (gRPC API)
- Реализован Payment Service (gRPC API)
- Интеграция через gRPC клиенты
- Thread-safe хранилища
- Graceful shutdown"
fi

# Проверка ветки
BRANCH=$(git branch --show-current 2>/dev/null || echo "main")
if [ -z "$BRANCH" ] || [ "$BRANCH" != "main" ]; then
    echo "🌿 Создание/переключение на ветку main..."
    git branch -M main
fi

echo ""
echo "✅ Готово к загрузке!"
echo ""
echo "📤 Для загрузки выполните:"
echo "   git push -u origin main"
echo ""
echo "💡 Если требуется аутентификация:"
echo "   1. Используйте Personal Access Token вместо пароля"
echo "   2. Или настройте SSH ключи для GitHub"
echo ""
read -p "Загрузить сейчас? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Загрузка в GitHub..."
    git push -u origin main
    echo "✅ Проект успешно загружен в GitHub!"
else
    echo "⏭️  Пропущено. Выполните 'git push -u origin main' позже."
fi

