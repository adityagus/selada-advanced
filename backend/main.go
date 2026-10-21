package main

import (
	"log"
	"os"

	"fmt"
	"go-api/config"
	"go-api/routes"
)

func loadDotEnv() {
	if err := config.LoadEnv(".env"); err != nil {
		log.Fatalf("Gagal membaca .env: %v", err)
	}
}

func main() {
	fmt.Println("Backend recompiled successfully!")
	loadDotEnv()

	// 1. Initialize Database Connections
	config.ConnectDatabases()

	// 2. Initialize router mapping
	r := routes.SetupRouter()

	// 3. Start Gin Server on port specified in env or fallback to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server berjalan di port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
