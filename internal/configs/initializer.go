package configs

import (
	"fmt"

	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/joho/godotenv"
)

func Init() {
	initEnv()
	initMigration()
}

func initEnv() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Errorf("failed to load the enviroments: %v", err))
	}
}

func initMigration() {
	database.ExecMigration()
}
