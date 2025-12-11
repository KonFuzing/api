package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import Internal Packages (Hexagonal Layers)
	"api/pkg/database" // Shared Database Package
	pb "api/proto"
	"api/services/arena/internal/adapters/client"
	"api/services/arena/internal/adapters/handler"
	"api/services/arena/internal/adapters/repository"
	"api/services/arena/internal/core/services"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	dsn := os.Getenv("DB_DSN")
	port := os.Getenv("ARENA_PORT")
	duelistTarget := os.Getenv("DUELIST_TARGET")

	// 2. Initialize Infrastructure (Database Singleton)
	db, err := database.GetInstance(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	fmt.Println("✅ Database connected successfully (Arena)")

	// 3. Initialize Infrastructure (gRPC Client Connection)
	// สร้าง Connection ไปยัง Duelist Service
	conn, err := grpc.NewClient(duelistTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Duelist Service: %v", err)
	}
	defer conn.Close() // ปิด Connection เมื่อจบโปรแกรม

	// สร้าง gRPC Client Instance จาก Proto
	grpcClient := pb.NewDuelistServiceClient(conn)

	// 4. Initialize Adapters (Secondary / Outbound)
	// 4.1 Repository Adapter (MySQL) สำหรับเก็บประวัติการต่อสู้
	repoAdapter := repository.NewMySQLRepository(db)

	// 4.2 Client Adapter (gRPC Wrapper) สำหรับดึงข้อมูล Cowboy
	clientAdapter := client.NewGrpcClientAdapter(grpcClient)

	// 5. Initialize Core Domain Service
	// Inject ทั้ง Client Adapter และ Repository Adapter เข้าไปใน Service Logic
	svc := services.NewArenaService(clientAdapter, repoAdapter)

	// 6. Initialize Primary Adapter (Inbound / HTTP Handler)
	httpHandler := handler.NewHttpHandler(svc)

	// 7. Register Routes & Start Server
	http.HandleFunc("/duel", httpHandler.HandleDuel)       // POST: เริ่มการต่อสู้
	http.HandleFunc("/history", httpHandler.HandleHistory) // GET: ดูประวัติ (มี Filter)

	fmt.Printf("⚔️  Arena Service (Hexagonal + HTTP) running on port :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
