#!/bin/bash
# Скрипт для удаления скомпилированных бинарников из git индекса

echo "🗑️  Удаляем скомпилированные файлы из git индекса..."

# Удаляем основные бинарники серверов
git rm --cached inventory/server 2>/dev/null || true
git rm --cached payment/server 2>/dev/null || true
git rm --cached order/server 2>/dev/null || true
git rm --cached server 2>/dev/null || true

# Удаляем возможные бинарники из других модулей
git rm --cached assembly/server 2>/dev/null || true
git rm --cached platform/server 2>/dev/null || true
git rm --cached iam/server 2>/dev/null || true
git rm --cached notification/server 2>/dev/null || true

# Удаляем тестовые бинарники
find . -name "*.test" -type f -exec git rm --cached {} \; 2>/dev/null || true

# Удаляем файлы покрытия
find . -name "*.out" -type f -exec git rm --cached {} \; 2>/dev/null || true

# Удаляем объектные файлы
find . -name "*.o" -type f -exec git rm --cached {} \; 2>/dev/null || true
find . -name "*.a" -type f -exec git rm --cached {} \; 2>/dev/null || true

# Удаляем временные файлы сборки
find . -name "*.tmp" -type f -exec git rm --cached {} \; 2>/dev/null || true
find . -name "*.bak" -type f -exec git rm --cached {} \; 2>/dev/null || true

echo "✅ Файлы удалены из git индекса (локальные файлы сохранены)"
echo "📝 Теперь эти файлы игнорируются через .gitignore"
echo ""
echo "Для применения изменений выполните:"
echo "  git add .gitignore"
echo "  git commit -m 'Remove compiled binaries and build artifacts from git tracking'"

