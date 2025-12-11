package database

import (
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ตัวแปร global (private) เก็บ instance
var (
	instance *gorm.DB
	once     sync.Once
	err      error
)

// GetInstance : ฟังก์ชันสำหรับเรียกใช้ DB (Singleton)
// จะทำการ connect แค่ครั้งแรกที่ถูกเรียก ครั้งต่อไปจะส่ง instance เดิมกลับไป
func GetInstance(dsn string) (*gorm.DB, error) {
	// sync.Once รับประกันว่า function ภายในจะทำงานแค่ 1 ครั้งตลอดอายุโปรแกรม
	once.Do(func() {
		log.Println("🔌 Initializing Database Connection (Singleton)...")

		instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return // ถ้า error ค่า err จะถูกเก็บไว้ return ออกไป
		}

		// (Optional) ตั้งค่า Connection Pool
		sqlDB, dbErr := instance.DB()
		if dbErr == nil {
			sqlDB.SetMaxIdleConns(10)  // จำนวน connection ที่เปิดรอไว้
			sqlDB.SetMaxOpenConns(100) // จำนวน connection สูงสุด
		}
	})

	return instance, err
}
