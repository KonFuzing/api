package main

import (
	pb "api/proto"
	"api/services/arena/internal/adapters/client"
	"api/services/arena/internal/adapters/handler"
	"api/services/arena/internal/adapters/repository"
	"api/services/arena/internal/core/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load("../../../.env")

	dbUrl := os.Getenv("DB_DSN")
	target := os.Getenv("DUELIST_TARGET")
	port := os.Getenv("ARENA_PORT")

	// 1. Init DB
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Init gRPC Connection
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := pb.NewDuelistServiceClient(conn)

	// 3. Init Adapters
	repoAdapter := repository.NewMySQLRepository(db)
	clientAdapter := client.NewGrpcClientAdapter(grpcClient)

	// 4. Init Service (Inject Adapters)
	svc := services.NewArenaService(clientAdapter, repoAdapter)

	// 5. Init HTTP Handler
	httpHandler := handler.NewHttpHandler(svc)

	// 6. Run
	http.HandleFunc("/duel", httpHandler.HandleDuel)
	http.HandleFunc("/history", httpHandler.HandleHistory)
	fmt.Printf("⚔️ Arena Service running on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
