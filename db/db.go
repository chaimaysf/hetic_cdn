package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	// "github.com/go-redis/redis/v8"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Chaimayousfi30!"
	DB_NAME     = "hetic_cdn"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

var DB *sql.DB
// var RedisClient *redis.Client
// var Ctx = context.Background() // ✅ Contexte Redis

func InitDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Erreur de connexion à PostgreSQL :", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Impossible de pinger PostgreSQL :", err)
	}

	fmt.Println("✅ Connexion PostgreSQL réussie !")
}

// Création des tables
func CreateTables() {
	query := `
	CREATE TABLE IF NOT EXISTS folders (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		parent_id INT REFERENCES folders(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		size BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		folder_id INT REFERENCES folders(id) ON DELETE CASCADE
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Erreur lors de la création des tables :", err)
	}
	fmt.Println("✅ Tables créées avec succès !")
}

