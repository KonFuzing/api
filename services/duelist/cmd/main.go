package main

import (
	"fmt"
	"log"
	"net"
	"os" // ยังต้องใช้ os เพื่อ override ชื่อ ENV เฉพาะของ service นี้

	"google.golang.org/grpc"

	// Import Packages
	"api/pkg/config" // ✅ เรียกใช้ Config Package
	"api/pkg/database"
	pb "api/proto"
	"api/services/duelist/internal/adapters/handler"
	"api/services/duelist/internal/adapters/repository"
	"api/services/duelist/internal/core/services"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// ⚠️ Override Port สำหรับ Duelist โดยเฉพาะ
	// (เพราะใน config กลางอาจจะเป็นค่า default)
	if p := os.Getenv("DUELIST_PORT"); p != "" {
		cfg.AppPort = p
	}

	// 2. Initialize Infrastructure (DB Singleton)
	// ใช้ค่าจาก cfg แทน os.Getenv
	db, err := database.GetInstance(cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// 3. Setup Layers (เหมือนเดิม)
	repoAdapter := repository.NewMySQLRepository(db)
	svc := services.NewDuelistService(repoAdapter)
	grpcHandler := handler.NewGrpcHandler(svc)

	// 4. Start Server (ใช้ Port จาก cfg)
	lis, err := net.Listen("tcp", ":"+cfg.AppPort)
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDuelistServiceServer(grpcServer, grpcHandler)

	fmt.Printf("🤠 Duelist Service running on port :%s\n", cfg.AppPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
