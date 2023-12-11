package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	// "io/ioutil"
	"mime/multipart"

	// "io/ioutil"
	"os"
)

// generateAESKey génère une clé AES
func generateAESKey() ([]byte, error) {
	key := make([]byte, 32) // Utilise une clé de 32 octets pour AES-256
	_, err := rand.Read(key)
	return key, err
}

// encryptFile chiffre le contenu d'un fichier avec AES
func encryptFile(key []byte, inputFile multipart.File, outputFile string) error {
	// Lire le contenu du fichier dans une variable
	plaintext, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	// Ajouter un padding pour s'assurer que la longueur du texte est un multiple du bloc AES
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	if padding != 0 {
		plaintext = append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return os.WriteFile(outputFile, ciphertext, 0644)
}

// decryptFile déchiffre le contenu d'un fichier chiffré avec AES
func decryptFile(key []byte, inputFile string) ([]byte, error) {
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}
