@echo off
setlocal enabledelayedexpansion

REM Определяем файл переменных окружения
if exist ".env" (
    set ENV_FILE=.env
) else if exist "env.txt" (
    set ENV_FILE=env.txt
) else (
    echo Файл .env или env.txt не найден
    exit /b 1
)

REM Загружаем переменные из файла
for /f "tokens=*" %%a in (%ENV_FILE%) do (
    set line=%%a
    if "!line:~0,1!" neq "#" (
        for /f "tokens=1,2 delims==" %%b in ("!line!") do (
            set %%b=%%c
        )
    )
)

REM Устанавливаем значения по умолчанию, если они не определены
if not defined DB_HOST set DB_HOST=localhost
if not defined DB_PORT set DB_PORT=5431
if not defined DB_USER set DB_USER=postgres
if not defined DB_PASSWORD set DB_PASSWORD=password
if not defined DB_NAME set DB_NAME=furni_swap

echo Применение миграции add_purchases.sql...

REM Выполняем SQL скрипт
set PGPASSWORD=%DB_PASSWORD%
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -f migrations/add_purchases.sql

if %ERRORLEVEL% neq 0 (
    echo Ошибка при применении миграции
    exit /b 1
)

echo Обновляем существующие объявления...

REM Устанавливаем статус available для всех существующих объявлений
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "UPDATE listings SET status = 'available' WHERE status IS NULL;"

if %ERRORLEVEL% neq 0 (
    echo Ошибка при обновлении статуса объявлений
    exit /b 1
)

echo Миграция успешно завершена
exit /b 0 