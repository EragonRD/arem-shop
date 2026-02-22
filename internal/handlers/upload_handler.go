package handlers

import (
	"errors"
	"net/http"

	"arem-shop/internal/services"
	"arem-shop/internal/utils"

	"github.com/gin-gonic/gin"
)

// UploadHandler s'occupe de recevoir les requêtes HTTP multipart pour les uploads de fichiers.
type UploadHandler struct {
	uploadService *services.UploadService
}

func NewUploadHandler(uploadService *services.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// UploadProductImage reçoit l'image depuis le frontend et la transfère au service.
func (h *UploadHandler) UploadProductImage(c *gin.Context) {
	// Récupérer le fichier via le champ de formulaire "image"
	file, err := c.FormFile("image")
	if err != nil {
		utils.JSONErrorWithCode(c, http.StatusBadRequest, "INVALID_FILE", "image field is required and must be a valid file")
		return
	}

	publicURL, err := h.uploadService.SaveImage(file)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFileTooLarge):
			utils.JSONErrorWithCode(c, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", err.Error())
		case errors.Is(err, services.ErrInvalidExt):
			utils.JSONErrorWithCode(c, http.StatusUnsupportedMediaType, "INVALID_EXTENSION", err.Error())
		default:
			utils.JSONErrorWithCode(c, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		}
		return
	}

	// Renvoyer l'URL publique générée en réponse JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url": publicURL,
		},
	})
}
