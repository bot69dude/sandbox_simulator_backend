package main

import (
	"log"

	"github.com/bot69dude/sandbox_simulator_backend.git/internal/config"
	"github.com/bot69dude/sandbox_simulator_backend.git/internal/repository"
)

func main() {
	config,err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := repository.InitDB()
	defer db.Close()

	log.Printf("Server starting on port %s in %s mode", config.ServerPort, config.Environment)

}
