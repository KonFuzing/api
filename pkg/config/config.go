package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBUrl         string
	DuelistTarget string // ใช้เฉพาะฝั่ง Arena
}

// LoadConfig : โหลดค่า Config ทั้งหมดทีเดียว
func LoadConfig() *Config {
	// พยายามโหลด .env แต่ถ้าไม่เจอก็ไม่เป็นไร (เผื่อรันบน Docker/Cloud)
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("⚠️  Note: .env file not found, using system env")
	}

	return &Config{
		// ฟังก์ชัน getEnv ช่วยเช็คค่า default ให้
		AppPort:       getEnv("APP_PORT", ""), // ใช้ชื่อกลางๆ เดี๋ยวไป override ใน main
		DBUrl:         getEnv("DB_DSN", ""),
		DuelistTarget: getEnv("DUELIST_TARGET", ""),
	}
}

// Helper function
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
