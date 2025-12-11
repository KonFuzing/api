package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// --- Config Part ---
type Config struct {
	Port string
}

// โหลด Config จาก Env ถ้าไม่มีให้ใช้ค่า Default
func LoadConfig() Config {
	// อ่านตัวแปรชื่อใหม่: DUELIST_PORT
	return Config{
		Port: getEnv("DUELIST_PORT", "8080"),
	}
}

// Helper function อ่านค่า Env
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// --- Domain & Handlers ---
type Cowboy struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Health   int     `json:"health"`
	Damage   int     `json:"damage"`
	Speed    int     `json:"speed"`
	Accuracy float64 `json:"accuracy"`
}

var (
	db = make(map[string]Cowboy)
	mu sync.Mutex
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: ไม่พบไฟล์ .env ที่ชั้นนอก (../.env)")
	}

	cfg := LoadConfig()

	http.HandleFunc("/cowboys", handleCreate)
	http.HandleFunc("/cowboys/", handleGetOne)

	fmt.Printf("🤠 Duelist Service running on port :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var c Cowboy
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	db[c.ID] = c
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
	fmt.Printf("[Duelist] Created: %s\n", c.Name)
}

func handleGetOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/cowboys/"):]
	mu.Lock()
	c, exists := db[id]
	mu.Unlock()
	if !exists {
		http.Error(w, "Cowboy not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(c)
}
