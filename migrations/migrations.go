package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/machariamuguku/golang-graphql-authentication/db"
	"github.com/machariamuguku/golang-graphql-authentication/models"
)

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// automigrate doesn't clear existing data
	// only adds none existent columns

	// automigrate one
	// db.AutoMigrate(models.GormUser{})

	// automigrate many
	// db.AutoMigrate(&GormUser{}, &GormProduct{}, &GormOrder{})

	db.DropTableIfExists(&models.GormUser{})
	db.CreateTable(&models.GormUser{})
}
