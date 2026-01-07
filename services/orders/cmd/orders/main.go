package main

import (
	"log"
	remna "orders/server/api/grpc/remna_client"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	addr := os.Getenv("REMNAVAWE_GRPC_BIND_URL")
	if len(addr) == 0 {
		panic("Bind url is required")
	}
	cfg := remna.Config{
		Address:        addr,
		RequestTimeout: 2 * time.Second,
		MaxRetries:     3,
	}
	client, err := remna.NewRemnaGRPCClient(cfg)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	user := client.GetUser("motbot")
	log.Println("UserResponse: ", user)
}
