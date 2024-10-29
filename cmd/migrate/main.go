package main

import (
	"log"
	"os"

	"github.com/ever864/ecommerce-psql/cmd/api"
	"github.com/ever864/ecommerce-psql/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(cmd string) error {
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		"postgres://efpl:postgres@localhost:5432/ecommerce?sslmode=disable",
	)
	if err != nil {
		return err
	}
	defer m.Close()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}

		log.Println("Migrations up successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return err
		}

		log.Println("Migrations down successfully")
	}
	return nil
}

func main() {
	// Verificar si hay argumentos de migración
	if len(os.Args) > 1 {
		cmd := os.Args[len(os.Args)-1]
		if cmd == "up" || cmd == "down" {
			if err := runMigrations(cmd); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	// Si no hay comandos de migración, inicia el servidor
	db.DatabaseConnection()
	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
