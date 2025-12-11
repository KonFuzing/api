package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import Packages
	"api/pkg/config" // ✅ เรียกใช้ Config Package
	"api/pkg/database"
	pb "api/proto"
	"api/services/arena/internal/adapters/client"
	"api/services/arena/internal/adapters/handler"
	"api/services/arena/internal/adapters/repository"
	"api/services/arena/internal/core/services"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// ⚠️ Override Port สำหรับ Arena
	if p := os.Getenv("ARENA_PORT"); p != "" {
		cfg.AppPort = p
	}

	// 2. Init DB (ใช้ cfg.DBUrl)
	db, err := database.GetInstance(cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// 3. Init gRPC Client (ใช้ cfg.DuelistTarget)
	conn, err := grpc.NewClient(cfg.DuelistTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Duelist: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewDuelistServiceClient(conn)

	// 4. Setup Layers (เหมือนเดิม)
	repoAdapter := repository.NewMySQLRepository(db)
	clientAdapter := client.NewGrpcClientAdapter(grpcClient)
	svc := services.NewArenaService(clientAdapter, repoAdapter)
	httpHandler := handler.NewHttpHandler(svc)

	// 5. Register Routes & Start
	http.HandleFunc("/duel", httpHandler.HandleDuel)
	http.HandleFunc("/history", httpHandler.HandleHistory)

	fmt.Printf("⚔️  Arena Service running on port :%s\n", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, nil); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
