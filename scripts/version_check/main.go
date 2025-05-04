package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Open("go.mod")
	if err != nil {
		log.Fatalf("Ошибка открытия файла go.mod: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	goVersionRegex := regexp.MustCompile(`^go \d+\.\d+(\.\d+)?$`)
	toolchainRegex := regexp.MustCompile(`^toolchain go\d+\.\d+(\.\d+)?$`)

	// Текущие стабильные версии
	stableGoVersion := "1.22.0"
	stableToolchain := "1.22.0"

	needsGoVersionUpdate := false
	needsToolchainUpdate := false

	// Считываем файл и находим строки с версиями
	for scanner.Scan() {
		line := scanner.Text()

		if goVersionRegex.MatchString(line) {
			// Заменяем версию на стабильную
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[0] == "go" {
				currentVersion := parts[1]
				if currentVersion != stableGoVersion {
					log.Printf("Найдена нестабильная версия Go: %s, заменяем на %s", currentVersion, stableGoVersion)
					lines = append(lines, fmt.Sprintf("go %s", stableGoVersion))
					needsGoVersionUpdate = true
					continue
				}
			}
		} else if toolchainRegex.MatchString(line) {
			// Заменяем версию toolchain на стабильную
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[0] == "toolchain" {
				currentToolchain := parts[1]
				if !strings.HasPrefix(currentToolchain, "go"+stableToolchain) {
					log.Printf("Найдена нестабильная версия toolchain: %s, заменяем на go%s", currentToolchain, stableToolchain)
					lines = append(lines, fmt.Sprintf("toolchain go%s", stableToolchain))
					needsToolchainUpdate = true
					continue
				}
			}
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения файла go.mod: %v", err)
	}

	if !needsGoVersionUpdate && !needsToolchainUpdate {
		log.Println("Версии Go в файле go.mod уже актуальны")
		return
	}

	// Записываем обновленный файл
	file, err = os.Create("go.mod.new")
	if err != nil {
		log.Fatalf("Ошибка создания временного файла: %v", err)
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	writer.Flush()
	file.Close()

	// Заменяем оригинальный файл
	err = os.Rename("go.mod.new", "go.mod")
	if err != nil {
		log.Fatalf("Ошибка замены файла go.mod: %v", err)
	}

	log.Println("Файл go.mod обновлен успешно")
}
