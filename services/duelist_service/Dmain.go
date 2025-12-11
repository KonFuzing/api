package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	// Import package ที่ generate มา (เปลี่ยน path ตาม project จริง)
	pb "github.com/yourusername/cowboy_arena/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// --- Config Part ---
type Config struct {
	Port string
}

func LoadConfig() Config {
	return Config{
		// gRPC มักไม่ใช้ 8080 (ที่เป็น HTTP) อาจจะใช้ 50051 หรือพอร์ตเดิมก็ได้
		Port: getEnv("DUELIST_PORT", "50051"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// --- Domain Implementation ---
// สร้าง struct เพื่อ implement interface ของ gRPC
type duelistServer struct {
	pb.UnimplementedDuelistServiceServer
	db map[string]*pb.CowboyResponse
	mu sync.Mutex
}

// Implement: CreateCowboy
func (s *duelistServer) CreateCowboy(ctx context.Context, req *pb.CreateCowboyRequest) (*pb.CowboyResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cowboy := &pb.CowboyResponse{
		Id:       req.Id,
		Name:     req.Name,
		Health:   req.Health,
		Damage:   req.Damage,
		Speed:    req.Speed,
		Accuracy: req.Accuracy,
	}

	s.db[req.Id] = cowboy
	fmt.Printf("[Duelist] Created: %s (ID: %s)\n", req.Name, req.Id)
	return cowboy, nil
}

// Implement: GetCowboy
func (s *duelistServer) GetCowboy(ctx context.Context, req *pb.GetCowboyRequest) (*pb.CowboyResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cowboy, exists := s.db[req.Id]
	if !exists {
		return nil, fmt.Errorf("cowboy not found with id: %s", req.Id)
	}

	return cowboy, nil
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: ไม่พบไฟล์ .env")
	}

	cfg := LoadConfig()

	// 1. สร้าง Listener บน TCP
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. สร้าง gRPC Server
	grpcServer := grpc.NewServer()

	// 3. Register Service ที่เรา implement เข้าไป
	service := &duelistServer{
		db: make(map[string]*pb.CowboyResponse),
	}
	pb.RegisterDuelistServiceServer(grpcServer, service)

	fmt.Printf("🤠 Duelist Service (gRPC) running on port :%s\n", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
