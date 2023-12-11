package main

import (
	"log"
	"net/http"
)

func setupHandlers() {
	// Configure les gestionnaires de route avec Gorilla Mux
	r := setupRoutes()

	// Utilise le routeur Gorilla Mux
	http.Handle("/", r)
}

func startServer() {
	// Démarre le serveur sur le port 8080
	port := ":8080"
	log.Printf("Serveur en cours d'écoute sur le port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	// Appelle la fonction de configuration de la base de données
	initDB()

	// Configure les gestionnaires de route avec Gorilla Mux
	setupHandlers()

	// Démarre le serveur HTTP
	startServer()

	log.Println("L'application a démarré.")
}
