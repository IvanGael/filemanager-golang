// routes.go

package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func setupRoutes() http.Handler {
	r := mux.NewRouter()

	// Routes de l'API
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/files", HandleGetFiles).Methods("GET")
	apiRouter.HandleFunc("/files/{fileName}", HandleReadFile).Methods("GET")
	apiRouter.HandleFunc("/files", HandleUploadFile).Methods("POST")
	apiRouter.HandleFunc("/files/{fileName}", HandleDeleteFile).Methods("DELETE")

	// Configurer CORS avec des options personnalisées
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	apiHandler := c.Handler(apiRouter)

	return apiHandler
}
