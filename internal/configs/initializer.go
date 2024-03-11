package configs

import (
	"fmt"

	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/joho/godotenv"
)

// Initialize all configurations
func Init() {
	initEnv()
	initMigration()
}

// Initialize enviroment variables
func initEnv() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Errorf("failed to load the enviroments: %v", err))
	}
}

// Initialize the migrations
func initMigration() {
	database.ExecMigration()
}
