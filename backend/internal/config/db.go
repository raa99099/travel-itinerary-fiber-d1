package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnvDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func ConnectDB() *gorm.DB {
	// Muat file .env jika ada (coba folder aktif lalu folder induk)
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../.env")
	}

	host := getEnvDefault("DB_HOST", "127.0.0.1")
	port := getEnvDefault("DB_PORT", "3306")
	user := getEnvDefault("DB_USER", "root")
	pass := os.Getenv("DB_PASSWORD")
	name := getEnvDefault("DB_NAME", "travel_itinerary_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	return db
}
