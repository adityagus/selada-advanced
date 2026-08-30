package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"fmt"
	"go-api/config"
	"go-api/routes"
)

func loadDotEnv() {
	file, err := os.Open(".env")
	fmt.Println("file : ", file)
	if err != nil {
		return // .env tidak ditemukan, biarkan membaca env sistem/docker
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			// Hanya set jika belum ada di env (misal dari docker compose)
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}

func main() {
	// 0. Load .env kustom untuk local development jika ada
	fmt.Print("testing")
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
