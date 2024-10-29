package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ever864/ecommerce-psql/cmd/api"
	"github.com/ever864/ecommerce-psql/db"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", os.Getenv("PORT")), db.DatabaseConnection())
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
