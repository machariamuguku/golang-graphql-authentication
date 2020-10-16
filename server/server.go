package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/joho/godotenv"
	golang_graphql_authentication "github.com/machariamuguku/golang-graphql-authentication"
	"github.com/machariamuguku/golang-graphql-authentication/db"
)

// load values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	port := os.Getenv("PORT")

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(golang_graphql_authentication.NewExecutableSchema(golang_graphql_authentication.Config{Resolvers: &golang_graphql_authentication.Resolver{DB: db}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
