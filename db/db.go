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
		Host     string
		User     string
		DBName   string
		Password string
		Port     int
		DBSource string
	}{
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		DBName:   os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Port:     port,
		DBSource: os.Getenv("DATABASE_SOURCE"),
	}

	log.Println("Connecting to database...")
	log.Println("Host:", config.Host)
	log.Println("User:", config.User)
	log.Println("Database:", config.DBName)
	log.Println("Port:", config.Port)

	// Construir string de conexión

	test := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	log.Println("Database connection string:", test)

	log.Println("Database connection string:", config.DBSource)

	// Abrir conexión

	db, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Verificar conexión

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal("Error pinging database: ", errPing)

	}

	log.Println("Successfully connected to database")

	return db
}
