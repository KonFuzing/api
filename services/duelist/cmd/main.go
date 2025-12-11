package main

import (
	pb "api/proto"
	"api/services/duelist/internal/adapters/handler"
	"api/services/duelist/internal/adapters/repository"
	"api/services/duelist/internal/core/services"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load("../../../.env") // Adjust path as needed

	dsn := os.Getenv("DB_DSN")
	port := os.Getenv("DUELIST_PORT")
	if port == "" {
		port = "50051"
	}

	// 1. Init Infrastructure (Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	// 2. Init Adapters (Repository)
	repo := repository.NewMySQLRepository(db)

	// 3. Init Core Service (Inject Repository)
	svc := services.NewDuelistService(repo)

	// 4. Init Primary Adapter (gRPC Handler)
	grpcHandler := handler.NewGrpcHandler(svc)

	// 5. Run Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterDuelistServiceServer(s, grpcHandler)

	fmt.Printf("🤠 Duelist Service running on :%s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
