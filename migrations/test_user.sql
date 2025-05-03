-- Создаем тестового пользователя с хешем пароля 'password123'
INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
VALUES ('test@example.com', '$2a$10$ZKqp8jEG9S.Di4VaKTobNeQbrK1GrPYWiAPZZYJEOK7jEuQgDGIG2', 'Тестовый', 'Пользователь', 'Москва', true, NOW())
ON CONFLICT (email) DO NOTHING;

-- Выводим ID созданного пользователя
SELECT id, email, name, last_name, is_verified FROM users WHERE email = 'test@example.com'; 