package handler

import (
	"io"
	"net/http"

	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
	"mediconnect/pkg/storage"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	blobService storage.BlobService
	authRepo    domain.AuthRepository
}

func NewUploadHandler(blobService storage.BlobService, authRepo domain.AuthRepository) *UploadHandler {
	return &UploadHandler{
		blobService: blobService,
		authRepo:    authRepo,
	}
}

// UploadKTP handles KTP image upload to Azure Blob Storage
func (h *UploadHandler) UploadKTP(c *gin.Context) {
	// Batasi ukuran file misal 5MB
	err := c.Request.ParseMultipartForm(5 << 20)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "File terlampau besar / bad request - "+err.Error())
		return
	}

	file, header, err := c.Request.FormFile("ktp")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gambar KTP wajib disertakan pada field 'ktp' - "+err.Error())
		return
	}
	defer file.Close()

	// Membaca isi file menjadi byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal membaca isi file - "+err.Error())
		return
	}

	// Content Type Header
	contentType := header.Header.Get("Content-Type")

	// Panggil service untuk upload ke blob
	if h.blobService == nil {
		response.Error(c, http.StatusInternalServerError, "Layanan penyimpanan awan sedang mengalami gangguan internal (Service Unavailable)")
		return
	}

	url, err := h.blobService.UploadKTP(c.Request.Context(), fileBytes, header.Filename, contentType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal melakukan proses upload ke Azure Blob Storage - "+err.Error())
		return
	}

	// Update the user's KtpURL in the database
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID type")
		return
	}

	if err := h.authRepo.UpdateUserKtpURL(c.Request.Context(), userIDStr, url); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user KTP URL")
		return
	}

	// Respon berhasil
	response.Success(c, http.StatusOK, "File KTP berhasil diupload dan dihubungkan ke profile user", map[string]interface{}{
		"url":               url,
		"filename_original": header.Filename,
		"size_bytes":        len(fileBytes),
		"content_type":      contentType,
	})
}
