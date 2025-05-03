-- Добавляем поле last_name в таблицу users
ALTER TABLE users ADD COLUMN last_name TEXT;

-- Обновляем существующих пользователей (устанавливаем пустую фамилию)
UPDATE users SET last_name = ''; 