-- SQL script to fix chat functionality issues
-- Run this script against your FurniSwap database

-- Update Russian names to English equivalents
UPDATE users 
SET name = 'Alexander' 
WHERE name = 'Александр';

UPDATE users 
SET name = 'Dmitry' 
WHERE name = 'Дмитрий';

UPDATE users 
SET name = 'Maxim' 
WHERE name = 'Максим';

UPDATE users 
SET name = 'Sergei' 
WHERE name = 'Сергей';

UPDATE users 
SET name = 'Ivan' 
WHERE name = 'Иван';

UPDATE users 
SET name = 'Andrey' 
WHERE name = 'Андрей';

UPDATE users 
SET name = 'Alexey' 
WHERE name = 'Алексей';

UPDATE users 
SET name = 'Artem' 
WHERE name = 'Артём';

UPDATE users 
SET name = 'Mikhail' 
WHERE name = 'Михаил';

UPDATE users 
SET name = 'Nikita' 
WHERE name = 'Никита';

UPDATE users 
SET name = 'Anna' 
WHERE name = 'Анна';

UPDATE users 
SET name = 'Maria' 
WHERE name = 'Мария';

UPDATE users 
SET name = 'Ekaterina' 
WHERE name = 'Екатерина';

UPDATE users 
SET name = 'Elena' 
WHERE name = 'Елена';

UPDATE users 
SET name = 'Olga' 
WHERE name = 'Ольга';

UPDATE users 
SET name = 'Natalia' 
WHERE name = 'Наталья';

UPDATE users 
SET name = 'Tatiana' 
WHERE name = 'Татьяна';

UPDATE users 
SET name = 'Yulia' 
WHERE name = 'Юлия';

UPDATE users 
SET name = 'Daria' 
WHERE name = 'Дарья';

UPDATE users 
SET name = 'Victoria' 
WHERE name = 'Виктория';

-- Update last names
UPDATE users 
SET last_name = 'Ivanov' 
WHERE last_name = 'Иванов';

UPDATE users 
SET last_name = 'Petrov' 
WHERE last_name = 'Петров';

UPDATE users 
SET last_name = 'Sidorov' 
WHERE last_name = 'Сидоров';

UPDATE users 
SET last_name = 'Smirnov' 
WHERE last_name = 'Смирнов';

UPDATE users 
SET last_name = 'Kuznetsov' 
WHERE last_name = 'Кузнецов';

UPDATE users 
SET last_name = 'Popov' 
WHERE last_name = 'Попов';

UPDATE users 
SET last_name = 'Sokolov' 
WHERE last_name = 'Соколов';

UPDATE users 
SET last_name = 'Mikhailov' 
WHERE last_name = 'Михайлов';

UPDATE users 
SET last_name = 'Novikov' 
WHERE last_name = 'Новиков';

UPDATE users 
SET last_name = 'Fedorov' 
WHERE last_name = 'Федоров';

UPDATE users 
SET last_name = 'Morozov' 
WHERE last_name = 'Морозов';

UPDATE users 
SET last_name = 'Volkov' 
WHERE last_name = 'Волков';

UPDATE users 
SET last_name = 'Alekseev' 
WHERE last_name = 'Алексеев';

UPDATE users 
SET last_name = 'Lebedev' 
WHERE last_name = 'Лебедев';

UPDATE users 
SET last_name = 'Semenov' 
WHERE last_name = 'Семенов';

UPDATE users 
SET last_name = 'Egorov' 
WHERE last_name = 'Егоров';

UPDATE users 
SET last_name = 'Pavlov' 
WHERE last_name = 'Павлов';

UPDATE users 
SET last_name = 'Kozlov' 
WHERE last_name = 'Козлов';

UPDATE users 
SET last_name = 'Stepanov' 
WHERE last_name = 'Степанов';

UPDATE users 
SET last_name = 'Nikolaev' 
WHERE last_name = 'Николаев';

-- Update the test user
UPDATE users
SET name = 'Test', last_name = 'User'
WHERE email = 'test@example.com';

-- Generate English email addresses for all users
UPDATE users
SET email = LOWER(CONCAT(name, '.', last_name, id, '@example.com'))
WHERE email NOT LIKE '%@example.com';

-- Check if query works
SELECT id, name, last_name, email FROM users LIMIT 10; 