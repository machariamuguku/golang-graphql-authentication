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

	db.DropTableIfExists(&models.GormUser{})
	db.CreateTable(&models.GormUser{})
}
