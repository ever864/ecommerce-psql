package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Driver de PostgreSQL
)

func DatabaseConnection() *sql.DB {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Obtener variables de entorno
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal("Error parsing database port:", err)
	}

	config := struct {
		Host   string
		User   string
		DBName string
		Port   int
	}{
		Host:   os.Getenv("DATABASE_HOST"),
		User:   os.Getenv("DATABASE_USER"),
		DBName: os.Getenv("DATABASE_NAME"),
		Port:   port,
	}

	// Construir string de conexión
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.DBName,
	)

	// Abrir conexión
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Successfully connected to database")

	return db
}
