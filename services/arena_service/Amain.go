package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// Import package ที่ generate มา
	pb "github.com/yourusername/cowboy_arena/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Config Part ---
type Config struct {
	Port          string
	DuelistSvcTarget string // เปลี่ยนจาก URL เป็น Target (host:port)
	DBUrl         string
}

func LoadConfig() Config {
	return Config{
		Port:          getEnv("ARENA_PORT", "8081"),
		// gRPC connect แบบ "host:port" ไม่ต้องมี http://
		DuelistSvcTarget: getEnv("DUELIST_TARGET", "localhost:50051"),
		DBUrl:         getEnv("DB_DSN", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var cfg Config
var db *gorm.DB
var grpcClient pb.DuelistServiceClient // เก็บ Client ไว้เรียกใช้

// --- Models ---
// เราสามารถใช้ struct เดิม หรือ map ข้อมูลจาก pb.CowboyResponse ก็ได้
// เพื่อความง่ายในการคำนวณขอใช้ struct เดิม แต่เขียน function แปลงข้อมูล
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

type BattleRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Fighter1ID string    `json:"fighter_1_id"`
	Fighter2ID string    `json:"fighter_2_id"`
	Winner     string    `json:"winner"`
	Logs       string    `gorm:"type:text" json:"logs"`
	CreatedAt  time.Time `json:"created_at"`
}

// --- Main ---

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: ไม่พบไฟล์ .env")
	}

	cfg = LoadConfig()

	// 1. เชื่อมต่อ Database
	initDB()

	// 2. เชื่อมต่อ gRPC ไปยัง Duelist Service
	// ใช้ WithTransportCredentials(insecure) สำหรับ local dev (ไม่มี SSL)
	conn, err := grpc.NewClient(cfg.DuelistSvcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Duelist Service: %v", err)
	}
	defer conn.Close()
	
	// สร้าง Client instance
	grpcClient = pb.NewDuelistServiceClient(conn)
	fmt.Printf("✅ Connected to Duelist Service at %s\n", cfg.DuelistSvcTarget)

	http.HandleFunc("/duel", handleDuel)
	http.HandleFunc("/history", handleHistory)

	fmt.Printf("⚔️ Arena Service (HTTP) running on port :%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, nil)
}

func initDB() {
	var err error
	db, err = gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&BattleRecord{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
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

	// Fetch ข้อมูล Cowboy ผ่าน gRPC
	c1, err1 := fetchCowboyGRPC(req.Fighter1ID)
	c2, err2 := fetchCowboyGRPC(req.Fighter2ID)

	if err1 != nil || err2 != nil {
		http.Error(w, fmt.Sprintf("Error fetching cowboys via gRPC: %v %v", err1, err2), http.StatusInternalServerError)
		return
	}

	result := simulateFight(c1, c2)

	record := BattleRecord{
		Fighter1ID: req.Fighter1ID,
		Fighter2ID: req.Fighter2ID,
		Winner:     result.Winner,
		Logs:       strings.Join(result.Logs, "\n"),
	}

	if err := db.Create(&record).Error; err != nil {
		log.Printf("⚠️ Error saving battle record: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	var history []BattleRecord
	result := db.Order("created_at desc").Find(&history)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// ฟังก์ชันใหม่: ดึงข้อมูลผ่าน gRPC
func fetchCowboyGRPC(id string) (*Cowboy, error) {
	// สร้าง Context พร้อม Timeout (Best Practice)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// เรียก gRPC method
	resp, err := grpcClient.GetCowboy(ctx, &pb.GetCowboyRequest{Id: id})
	if err != nil {
		return nil, err
	}

	// แปลงจาก gRPC struct กลับมาเป็น Domain struct (Cowboy)
	return &Cowboy{
		ID:       resp.Id,
		Name:     resp.Name,
		Health:   int(resp.Health),
		Damage:   int(resp.Damage),
		Speed:    int(resp.Speed),
		Accuracy: resp.Accuracy,
	}, nil
}

func simulateFight(c1, c2 *Cowboy) DuelResult {
	var logs []string
	logs = append(logs, fmt.Sprintf("Match Start: %s VS %s", c1.Name, c2.Name))
	
	// Dummy Logic เดิม
	winner := c1.Name
	logs = append(logs, fmt.Sprintf("%s fires a shot with %d damage!", c1.Name, c1.Damage))
	logs = append(logs, fmt.Sprintf("%s wins!", winner))

	return DuelResult{Winner: winner, Logs: logs}
}