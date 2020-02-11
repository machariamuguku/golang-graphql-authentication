package main

import (
	"github.com/machariamuguku/golang-graphql-authentication/db"
	"github.com/machariamuguku/golang-graphql-authentication/models"
)

func main() {
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
