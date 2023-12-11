//db.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	// Charge les variables d'environnement à partir du fichier .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Impossible de charger le fichier .env")
	}

	// Initialise la connexion à la base de données
	setupDB()
}

func setupDB() {
	var err error

	// Configuration de la connexion à PostgreSQL avec Gorm
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Vérifie si la base de données existe, sinon la crée
	if err = createDatabaseIfNotExists(dsn); err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate pour créer automatiquement les tables
	err = db.AutoMigrate(&File{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connexion à la base de données réussie!")
}

// createDatabaseIfNotExists crée la base de données si elle n'existe pas
func createDatabaseIfNotExists(dsn string) error {
	// Ouvre une connexion temporaire à PostgreSQL pour exécuter la requête de création de base de données
	tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Obtient l'objet *sql.DB sous-jacent à partir de gorm.DB
	sqlDB, err := tempDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Vérifie si la base de données existe
	result := tempDB.Exec("SELECT 1 FROM pg_database WHERE datname = ?", os.Getenv("DB_NAME"))
	if result.Error != nil {
		return result.Error
	}

	// Si la base de données n'existe pas, crée-la
	if result.RowsAffected == 0 {
		// Utiliser un superutilisateur (superuser) pour créer la base de données
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s WITH OWNER = postgres", os.Getenv("DB_NAME"))
		err := tempDB.Exec(createDBQuery).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// CloseDB ferme la connexion à la base de données
func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	log.Println("Fermeture de la connexion à la base de données.")
}
