// file.services.go

package main

import (
	"encoding/base64"
	"encoding/json"
	"time"

	// "io"

	// "io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Endpoint pour télécharger un fichier
// func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Impossible de récupérer le fichier", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Créer un fichier sur le serveur
// 	filePath := "./uploads/" + handler.Filename
// 	newFile, err := os.Create(filePath)
// 	if err != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Impossible de créer le fichier sur le serveur", http.StatusInternalServerError)
// 		return
// 	}
// 	defer newFile.Close()

// 	// Copier le contenu du fichier téléchargé vers le nouveau fichier
// 	_, err = io.Copy(newFile, file)
// 	if err != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Impossible de copier le contenu du fichier", http.StatusInternalServerError)
// 		return
// 	}

// 	// Enregistre les informations du fichier dans la base de données
// 	db.Create(&File{FileName: handler.Filename, FilePath: filePath})

// 	// Réponse réussie
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Fichier téléchargé avec succès"))
// }

// Endpoint pour lire un fichier
// func HandleReadFile(w http.ResponseWriter, r *http.Request) {
// 	fileName := mux.Vars(r)["fileName"]

// 	var file File
// 	result := db.Where("file_name = ?", fileName).First(&file)
// 	if result.Error != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
// 		return
// 	}

// 	// Ouvrir et lire le contenu du fichier
// 	fileContent, err := os.ReadFile(file.FilePath)
// 	if err != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Impossible de lire le contenu du fichier", http.StatusInternalServerError)
// 		return
// 	}

// 	// Envoyer le contenu du fichier en tant que réponse
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(fileContent)
// }

var encryptKey, _ = generateAESKey()

// Fonction pour uploader un fichier
func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		// Gérer l'erreur
		http.Error(w, "Impossible de récupérer le fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	// Chiffrer le contenu du fichier avant de l'enregistrer
	encryptedFilePath := "./uploads/" + fileName
	err = encryptFile(encryptKey, file, encryptedFilePath)
	if err != nil {
		// Gérer l'erreur
		http.Error(w, "Impossible de chiffrer le fichier", http.StatusInternalServerError)
		return
	}

	// Enregistre les informations du fichier dans la base de données
	db.Create(&File{FileName: fileName, FilePath: encryptedFilePath})

	// Réponse réussie
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fichier téléchargé avec succès"))
}

// Fonction pour lire un fichier
// func HandleReadFile(w http.ResponseWriter, r *http.Request) {
// 	fileName := mux.Vars(r)["fileName"]

// 	var file File
// 	result := db.Where("file_name = ?", fileName).First(&file)
// 	if result.Error != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
// 		return
// 	}

// 	// Déchiffrer le fichier avant de l'envoyer en réponse
// 	decryptedContent, err := decryptFile(encryptKey, file.FilePath)
// 	if err != nil {
// 		// Gérer l'erreur
// 		http.Error(w, "Impossible de déchiffrer le fichier", http.StatusInternalServerError)
// 		return
// 	}

// 	// Envoyer le contenu du fichier en tant que réponse
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(decryptedContent)
// }

func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]

	var file File
	result := db.Where("file_name = ?", fileName).First(&file)
	if result.Error != nil {
		// Gérer l'erreur
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	// Déchiffrer le fichier avant de l'envoyer en réponse
	decryptedContent, err := decryptFile(encryptKey, file.FilePath)
	if err != nil {
		// Gérer l'erreur
		http.Error(w, "Impossible de déchiffrer le fichier", http.StatusInternalServerError)
		return
	}

	// Convertir le contenu du fichier en base64
	base64Content := base64.StdEncoding.EncodeToString(decryptedContent)

	// Envoyer le contenu du fichier en tant que réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"content": base64Content})
}

// Générer un nom de fichier unique en ajoutant un suffixe numérique
func generateUniqueFileName(baseName string) string {
	currentTime := time.Now().String()
	return "File" + currentTime
}

// Fonction pour récupérer tous les fichiers
func HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	var Files []File
	if err := db.Find(&Files).Error; err != nil {
		http.Error(w, "Erreur lors de la récupération des fichiers", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(Files)
	if err != nil {
		http.Error(w, "Erreur lors de la sérialisation des fichiers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Fonction pour supprimer un fichier
func HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]

	var file File
	result := db.Where("file_name = ?", fileName).First(&file)
	if result.Error != nil {
		// Gérer l'erreur
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	// Supprimer le fichier du système de fichiers
	err := os.Remove(file.FilePath)
	if err != nil {
		// Gérer l'erreur
		http.Error(w, "Impossible de supprimer le fichier du système de fichiers", http.StatusInternalServerError)
		return
	}

	// Supprimer le fichier de la base de données
	db.Delete(&file)

	// Réponse réussie
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fichier supprimé avec succès"))
}
