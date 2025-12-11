package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// Import Internal Packages (Hexagonal Layers)
	"api/pkg/database" // Shared Database Package
	pb "api/proto"
	"api/services/duelist/internal/adapters/handler"
	"api/services/duelist/internal/adapters/repository"
	"api/services/duelist/internal/core/services"
)

func main() {
	// 1. Load Environment Variables
	// ปรับ Path .env ตามความเหมาะสมของโครงสร้าง Folder จริง
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}

	dsn := os.Getenv("DB_DSN")
	port := os.Getenv("DUELIST_PORT")

	// 2. Initialize Infrastructure (Database Singleton)
	db, err := database.GetInstance(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	fmt.Println("✅ Database connected successfully (Duelist)")

	// 3. Initialize Adapters (Secondary / Outbound)
	// Inject DB instance เข้าไปใน Repository
	repoAdapter := repository.NewMySQLRepository(db)

	// 4. Initialize Core Domain Service
	// Inject Repository เข้าไปใน Service (Business Logic)
	svc := services.NewDuelistService(repoAdapter)

	// 5. Initialize Primary Adapter (Inbound / Handler)
	// Inject Service เข้าไปใน gRPC Handler
	grpcHandler := handler.NewGrpcHandler(svc)

	// 6. Start gRPC Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDuelistServiceServer(grpcServer, grpcHandler)

	fmt.Printf("🤠 Duelist Service (Hexagonal + gRPC) running on port :%s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve gRPC: %v", err)
	}
}
