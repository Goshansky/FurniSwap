#!/bin/bash

# Получаем параметры подключения из env.txt или .env
if [ -f ".env" ]; then
  source .env
elif [ -f "env.txt" ]; then
  source env.txt
else
  echo "Файл .env или env.txt не найден"
  exit 1
fi

# Устанавливаем значения по умолчанию, если они не определены
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5431}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"
DB_NAME="${DB_NAME:-furni_swap}"

echo "Применение миграции add_purchases.sql..."

# Выполняем SQL скрипт
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f migrations/add_purchases.sql

if [ $? -eq 0 ]; then
  echo "Миграция успешно применена"
else
  echo "Ошибка при применении миграции"
  exit 1
fi

echo "Обновляем существующие объявления..."

# Устанавливаем статус available для всех существующих объявлений, если это ещё не было сделано
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "UPDATE listings SET status = 'available' WHERE status IS NULL;"

if [ $? -eq 0 ]; then
  echo "Статус объявлений успешно обновлен"
else
  echo "Ошибка при обновлении статуса объявлений"
  exit 1
fi

echo "Миграция успешно завершена" 