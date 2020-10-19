package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// gorm postgres dialect- comment to justify underscore import
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB *gorm.DB
type DB struct {
	*gorm.DB
}

// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {

	// database variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	//Build connection string
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)

	// connect to the db
	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		log.Printf("ConnectDB: %v", err)
	} else {
		log.Printf("ConnectDB: successfully connected to the %v database", dbName)
	}

	return &DB{db}, nil
}
