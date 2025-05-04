package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.Println("Запуск всех скриптов обслуживания...")

	// Получаем текущий рабочий каталог
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка получения текущего каталога: %v", err)
	}

	// Запускаем скрипт сброса и заполнения БД
	fixdbPath := filepath.Join(wd, "scripts", "new_tools", "fixdb")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = fixdbPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Запуск скрипта сброса и заполнения БД...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Ошибка при запуске скрипта сброса и заполнения БД: %v", err)
	}

	// Запускаем скрипт исправления email-адресов
	emailsPath := filepath.Join(wd, "scripts", "new_tools", "emails")
	cmd = exec.Command("go", "run", "main.go")
	cmd.Dir = emailsPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Запуск скрипта исправления email-адресов...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Ошибка при запуске скрипта исправления email-адресов: %v", err)
	}

	fmt.Println("Все скрипты выполнены успешно!")
}
