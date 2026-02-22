package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var (
	ErrFileTooLarge   = errors.New("file is too large (max 5MB)")
	ErrInvalidExt     = errors.New("invalid file extension (only jpg, png, webp allowed)")
	ErrUploadFailed   = errors.New("failed to save the uploaded file")
)

// UploadService gère l'enregistrement local des fichiers uploadés par les utilisateurs
// (ex: images de produits).
type UploadService struct {
	UploadDirectory string
	BaseURL         string
}

// NewUploadService construit un service d'upload.
// uploadDirectory est le chemin local du dossier où seront stockés les fichiers (par exemple "./uploads").
// baseURL est l'URL publique de base de l'API (ex: "http://localhost:8080").
func NewUploadService(uploadDirectory, baseURL string) *UploadService {
	// Créer le dossier s'il n'existe pas déjà
	_ = os.MkdirAll(uploadDirectory, 0755)

	return &UploadService{
		UploadDirectory: uploadDirectory,
		BaseURL:         baseURL,
	}
}

// SaveImage prend un multipart.FileHeader et l'enregistre sur le disque dur
// avec un nom unique (UUID) pour éviter les collisions écrasant d'autres fichiers.
// Elle retourne l'URL absolue publique pour y accéder.
func (s *UploadService) SaveImage(fileHeader *multipart.FileHeader) (string, error) {
	// 1Vérification de la taille (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		return "", ErrFileTooLarge
	}

	// 2. Vérification de l'extension
	ext := filepath.Ext(fileHeader.Filename)
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !validExts[ext] {
		return "", ErrInvalidExt
	}

	// 3. Génération d'un nom de fichier unique sécurisé avec un UUID
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := filepath.Join(s.UploadDirectory, newFilename)

	// 4. Ouverture du fichier entrant
	src, err := fileHeader.Open()
	if err != nil {
		return "", ErrUploadFailed
	}
	defer src.Close()

	// 5. Création du fichier de destination sur le serveur
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", ErrUploadFailed
	}
	defer dst.Close()

	// 6. Copie des octets (stream) pour sauvegarder le fichier
	if _, err = io.Copy(dst, src); err != nil {
		return "", ErrUploadFailed
	}

	// 7. Retourner l'URL absolue (Ex: http://localhost:8080/uploads/uuid.jpg)
	publicURL := fmt.Sprintf("%s/uploads/%s", s.BaseURL, newFilename)
	return publicURL, nil
}
