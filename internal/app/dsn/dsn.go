package dsn

import (
	"fmt"
	"os"
)

func FromEnv() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	// ✅ ПРАВИЛЬНАЯ строка подключения для GORM + PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, pass, name, port,
	)

	return dsn
}
