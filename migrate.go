package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
	migrationDir = "migrations" // путь к директории с миграциями
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Получаем параметры подключения к базе данных
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("SSL_MODE")

	// Строка подключения к базе данных
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	// Указываем флаги для команд (up, down, force, version, create)
	action := flag.String("action", "", "Migration action: up, down, force, version, create")
	migrationName := flag.String("name", "", "Migration name (only for 'create')")
	steps := flag.Int("steps", 0, "Number of steps to migrate (only for 'down')")
	flag.Parse()

	// Проверка, что указано действие
	if *action == "" {
		fmt.Println("Usage: -action=[up|down|force|version|create] [-steps=N] [-name=migration_name]")
		os.Exit(1)
	}

	// Создаем миграции, если команда "create"
	if *action == "create" && *migrationName != "" {
		createMigration(*migrationName)
		os.Exit(0)
	}

	// Создаем мигратор
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationDir),
		dsn,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// Выполнение действий в зависимости от флага
	switch *action {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration up applied successfully")

	case "down":
		if *steps > 0 {
			err = m.Steps(-(*steps))
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Migration down applied successfully")

	case "force":
		err = m.Force(*steps)
		if err != nil {
			log.Fatalf("Migration force failed: %v", err)
		}
		fmt.Println("Migration force applied successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Migration version failed: %v", err)
		}
		fmt.Printf("Current migration version: %d, dirty: %v\n", version, dirty)

	default:
		fmt.Println("Unknown action. Available actions: up, down, force, version, create")
		os.Exit(1)
	}
}

// createMigration создает миграционные файлы с заданным именем
func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	upFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	// Создаем файл .up.sql
	if err := createFile(upFile); err != nil {
		log.Fatalf("Could not create migration up file: %v", err)
	}
	// Создаем файл .down.sql
	if err := createFile(downFile); err != nil {
		log.Fatalf("Could not create migration down file: %v", err)
	}

	fmt.Printf("Created migration: \n%s\n%s\n", upFile, downFile)
}

// createFile создает пустой файл для миграции
func createFile(filePath string) error {
	_, err := os.Stat(migrationDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(migrationDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create migrations directory: %v", err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
