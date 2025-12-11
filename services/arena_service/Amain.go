package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings" // ใช้สำหรับรวม Logs เป็น string เดียว
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Config Part ---
type Config struct {
	Port          string
	DuelistSvcURL string
	DBUrl         string // เพิ่มตัวแปรรับ DSN
}

func LoadConfig() Config {
	return Config{
		Port:          getEnv("ARENA_PORT", "8081"),
		DuelistSvcURL: getEnv("DUELIST_URL", "http://localhost:8080"),
		DBUrl:         getEnv("DB_DSN", ""), // อ่านค่า Connection String
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var cfg Config
var db *gorm.DB // Global variable สำหรับ Database

// --- Models ---
// Cowboy (เหมือนเดิม)
type Cowboy struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Health   int     `json:"health"`
	Damage   int     `json:"damage"`
	Speed    int     `json:"speed"`
	Accuracy float64 `json:"accuracy"`
}

type DuelRequest struct {
	Fighter1ID string `json:"fighter_1"`
	Fighter2ID string `json:"fighter_2"`
}

type DuelResult struct {
	Winner string   `json:"winner"`
	Logs   []string `json:"battle_logs"`
}

// --- Database Models ---
// BattleRecord: ตารางเก็บประวัติการดวล
type BattleRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Fighter1ID string    `json:"fighter_1_id"`
	Fighter2ID string    `json:"fighter_2_id"`
	Winner     string    `json:"winner"`
	Logs       string    `gorm:"type:text" json:"logs"` // เก็บ Log ยาวๆ เป็น Text
	CreatedAt  time.Time `json:"created_at"`
}

// --- Main ---

func main() {
	// Load Env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: ไม่พบไฟล์ .env ที่ชั้นนอก")
	}

	cfg = LoadConfig()

	// 1. เชื่อมต่อ Database
	initDB()

	http.HandleFunc("/duel", handleDuel)
	// แถม: API ดูประวัติการดวล
	http.HandleFunc("/history", handleHistory)

	fmt.Printf("⚔️ Arena Service running on port :%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, nil)
}

// ฟังก์ชันเชื่อมต่อ Database
func initDB() {
	var err error
	// เปิด Connection ไปยัง Postgres
	db, err = gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Auto Migrate: สร้างตาราง BattleRecord อัตโนมัติถ้ายังไม่มี
	err = db.AutoMigrate(&BattleRecord{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	fmt.Println("✅ Database connected and migrated successfully!")
}

func handleDuel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req DuelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch ข้อมูล Cowboy
	c1, err1 := fetchCowboy(req.Fighter1ID)
	c2, err2 := fetchCowboy(req.Fighter2ID)

	if err1 != nil || err2 != nil {
		http.Error(w, fmt.Sprintf("Error fetching cowboys: %v %v", err1, err2), http.StatusInternalServerError)
		return
	}

	// คำนวณผลการต่อสู้
	result := simulateFight(c1, c2)

	// 2. บันทึกผลลง Database
	record := BattleRecord{
		Fighter1ID: req.Fighter1ID,
		Fighter2ID: req.Fighter2ID,
		Winner:     result.Winner,
		// แปลง []string เป็น string ยาวๆ คั่นด้วยบรรทัดใหม่ เพื่อเก็บลง DB
		Logs: strings.Join(result.Logs, "\n"),
	}

	// สั่ง Save ลง DB (Async เพื่อไม่ให้ Response ช้าเกินไป หรือจะรอผลก็ได้)
	if err := db.Create(&record).Error; err != nil {
		log.Printf("⚠️ Error saving battle record: %v", err)
	} else {
		log.Printf("✅ Battle record saved! ID: %d", record.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// API สำหรับดูประวัติการดวลทั้งหมด
func handleHistory(w http.ResponseWriter, r *http.Request) {
	var history []BattleRecord
	// ดึงข้อมูลทั้งหมดจากตาราง battle_records เรียงตามเวลาล่าสุด
	result := db.Order("created_at desc").Find(&history)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// fetchCowboy (เหมือนเดิม)
func fetchCowboy(id string) (*Cowboy, error) {
	url := fmt.Sprintf("%s/cowboys/%s", cfg.DuelistSvcURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cowboy not found")
	}

	var c Cowboy
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func simulateFight(c1, c2 *Cowboy) DuelResult {
	// ... Logic การต่อสู้ (สมมติว่าเหมือนเดิม) ...
	var logs []string
	logs = append(logs, fmt.Sprintf("Match Start: %s VS %s", c1.Name, c2.Name))

	// Dummy Logic
	winner := c1.Name
	logs = append(logs, fmt.Sprintf("%s fires a shot!", c1.Name))
	logs = append(logs, fmt.Sprintf("%s wins!", winner))

	return DuelResult{Winner: winner, Logs: logs}
}
