package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlexWilliam12/silent-signal/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Execute migrations to database
func ExecMigration() {
	db := OpenConn()
	if err := db.AutoMigrate(&models.User{}, &models.PrivateMessage{}, &models.Group{}, &models.GroupMessage{}); err != nil {
		panic(fmt.Errorf("failed to execute migration: %v", err))
	}
}

// Open a database connection
func OpenConn() *gorm.DB {

	// Connect on the database
	db, err := gorm.Open(getConn(), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}

	// Get the database connection
	conn, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get database connection: %v", err))
	}

	// Convert string to int
	idle, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNECTIONS"))
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	// Convert string to int
	max, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	// Max idle connections
	conn.SetMaxIdleConns(idle)

	// Max openned connections
	conn.SetMaxOpenConns(max)

	return db
}

// Get the connection to database
func getConn() gorm.Dialector {
	return postgres.Open(fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s
	sslmode=disable
	TimeZone=America/Sao_Paulo`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT")))
}
