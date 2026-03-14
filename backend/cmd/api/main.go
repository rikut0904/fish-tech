package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"fish-tech/internal/infrastructure/router"
	"fish-tech/internal/shared/timeutil"
)

func main() {
	time.Local = timeutil.JSTLocation()
	loadEnvFiles(".env", "backend/.env")

	e, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}

// loadEnvFiles は既存の環境変数を優先しつつ .env を読み込みます。
func loadEnvFiles(paths ...string) {
	for _, path := range paths {
		if err := loadEnvFile(path); err == nil {
			return
		}
	}
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
