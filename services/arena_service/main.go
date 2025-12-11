package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// --- Config Part ---
type Config struct {
	Port          string
	DuelistSvcURL string // เก็บ URL ของ Service ปลายทาง
}

func LoadConfig() Config {
	// อ่านตัวแปรชื่อใหม่: ARENA_PORT และ DUELIST_URL
	return Config{
		Port:          getEnv("ARENA_PORT", "8081"),
		DuelistSvcURL: getEnv("DUELIST_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Global Config variable
var cfg Config

// --- Models ---
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

// --- Main ---

func main() {

	// Point ไปที่ไฟล์ .env ชั้นนอกสุด (../.env)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: ไม่พบไฟล์ .env ที่ชั้นนอก (../.env)")
	}

	cfg = LoadConfig()

	http.HandleFunc("/duel", handleDuel)

	fmt.Printf("⚔️ Arena Service running on port :%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, nil)
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

	c1, err1 := fetchCowboy(req.Fighter1ID)
	c2, err2 := fetchCowboy(req.Fighter2ID)

	if err1 != nil || err2 != nil {
		http.Error(w, fmt.Sprintf("Error fetching cowboys: %v %v", err1, err2), http.StatusInternalServerError)
		return
	}

	result := simulateFight(c1, c2)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// fetchCowboy ใช้ URL จาก Config แทน Hardcode
func fetchCowboy(id string) (*Cowboy, error) {
	// ใช้ cfg.DuelistSvcURL ที่โหลดไว้
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
	// ... (Logic การต่อสู้เหมือนเดิม ไม่เปลี่ยนแปลง) ...
	// เพื่อความกระชับ ขอละไว้ในฐานที่เข้าใจครับ
	// (สามารถ copy logic เดิมมาใส่ได้เลย)

	var logs []string
	logs = append(logs, fmt.Sprintf("Match Start: %s VS %s", c1.Name, c2.Name))
	winner := c1.Name
	// Dummy Result เพื่อให้ Code compile ผ่านในตัวอย่างนี้
	return DuelResult{Winner: winner, Logs: logs}
}
