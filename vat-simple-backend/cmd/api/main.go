package main

import (
	"fmt"
	"log"

	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/api/router"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/config"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/database"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils" // TYPO FIXED HERE
)

func main() {
	// 1. Load configuration from .env file
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	log.Printf("Configuration loaded successfully. Server will run on port: %s", cfg.ServerPort)

	// 2. Initialize Database
	db, err := database.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.CloseMySQL(db)
	log.Println("Successfully connected to MySQL database!")

	// 3. Initialize Utilities
	utils.InitValidator()
	utils.InitJWT(cfg)

	// 4. Setup Router
	r := router.SetupRouter(db)

	// 5. Start Server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is listening on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
