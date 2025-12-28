package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ChunkSize     = 5 * 1024 * 1024 // 5MB chunks
	TempUploadDir = "uploads/temp"
	MaxChunkAge   = 24 * time.Hour // Cleanup old chunks after 24h
)

type ChunkedUploadHandler struct {
	DB *sql.DB
}

func NewChunkedUploadHandler(db *sql.DB) *ChunkedUploadHandler {
	// Create temp directory if not exists
	os.MkdirAll(TempUploadDir, 0755)
	return &ChunkedUploadHandler{DB: db}
}

// InitiateUpload creates a new upload session
func (h *ChunkedUploadHandler) InitiateUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get file info from request
	fileName := r.FormValue("file_name")
	fileSizeStr := r.FormValue("file_size")
	_ = r.FormValue("media_type") // Reserved for future use
	_ = r.FormValue("mime_type")  // Reserved for future use

	if fileName == "" || fileSizeStr == "" {
		sendErrorResponse(w, "Missing file_name or file_size", http.StatusBadRequest)
		return
	}

	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Invalid file_size", http.StatusBadRequest)
		return
	}

	// Generate upload ID
	uploadID := uuid.New().String()

	// Create temp directory for this upload
	uploadDir := filepath.Join(TempUploadDir, uploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		sendErrorResponse(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Calculate total chunks
	totalChunks := int((fileSize + ChunkSize - 1) / ChunkSize)

	fmt.Printf("📤 [CHUNKED] Инициализация загрузки: upload_id=%s, user_id=%d, file=%s, size=%d, chunks=%d\n",
		uploadID, userID, fileName, fileSize, totalChunks)

	sendSuccessResponse(w, map[string]interface{}{
		"upload_id":    uploadID,
		"chunk_size":   ChunkSize,
		"total_chunks": totalChunks,
	})
}

// UploadChunk handles individual chunk upload
func (h *ChunkedUploadHandler) UploadChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get chunk info
	uploadID := r.FormValue("upload_id")
	chunkIndexStr := r.FormValue("chunk_index")

	if uploadID == "" || chunkIndexStr == "" {
		sendErrorResponse(w, "Missing upload_id or chunk_index", http.StatusBadRequest)
		return
	}

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		sendErrorResponse(w, "Invalid chunk_index", http.StatusBadRequest)
		return
	}

	// Get chunk data
	file, _, err := r.FormFile("chunk")
	if err != nil {
		sendErrorResponse(w, "Failed to read chunk", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save chunk to temp directory
	uploadDir := filepath.Join(TempUploadDir, uploadID)
	chunkPath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", chunkIndex))

	dst, err := os.Create(chunkPath)
	if err != nil {
		sendErrorResponse(w, "Failed to save chunk", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	chunkSize, err := io.Copy(dst, file)
	if err != nil {
		sendErrorResponse(w, "Failed to write chunk", http.StatusInternalServerError)
		return
	}

	fmt.Printf("📦 [CHUNKED] Chunk загружен: upload_id=%s, chunk=%d, size=%d\n", uploadID, chunkIndex, chunkSize)

	sendSuccessResponse(w, map[string]interface{}{
		"chunk_index": chunkIndex,
		"uploaded":    true,
	})
}

// CompleteUpload assembles chunks and processes the file
func (h *ChunkedUploadHandler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get upload info
	uploadID := r.FormValue("upload_id")
	fileName := r.FormValue("file_name")
	mediaType := r.FormValue("media_type")
	mimeType := r.FormValue("mime_type")
	totalChunksStr := r.FormValue("total_chunks")

	if uploadID == "" || fileName == "" {
		sendErrorResponse(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		sendErrorResponse(w, "Invalid total_chunks", http.StatusBadRequest)
		return
	}

	fmt.Printf("🔗 [CHUNKED] Сборка файла: upload_id=%s, chunks=%d\n", uploadID, totalChunks)

	// Assemble chunks
	uploadDir := filepath.Join(TempUploadDir, uploadID)

	// Generate final file path
	ext := filepath.Ext(fileName)
	finalFileName := uuid.New().String() + ext
	now := time.Now()
	relativePath := filepath.Join("users", strconv.Itoa(userID), mediaType+"s",
		strconv.Itoa(now.Year()), fmt.Sprintf("%02d", now.Month()), finalFileName)
	fullPath := filepath.Join(UploadDir, relativePath)

	// Create directory
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		sendErrorResponse(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// Create final file
	finalFile, err := os.Create(fullPath)
	if err != nil {
		sendErrorResponse(w, "Failed to create final file", http.StatusInternalServerError)
		return
	}
	defer finalFile.Close()

	// Assemble chunks
	var totalSize int64
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			sendErrorResponse(w, fmt.Sprintf("Missing chunk %d", i), http.StatusBadRequest)
			os.Remove(fullPath)
			return
		}

		size, err := io.Copy(finalFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			sendErrorResponse(w, "Failed to assemble chunks", http.StatusInternalServerError)
			os.Remove(fullPath)
			return
		}
		totalSize += size
	}

	fmt.Printf("✅ [CHUNKED] Файл собран: %s, size=%d\n", fullPath, totalSize)

	// Save to database FIRST (with original file)
	query := `
		INSERT INTO user_media (user_id, file_name, original_name, file_path, file_size, mime_type, media_type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := h.DB.Exec(query, userID, finalFileName, fileName, relativePath, totalSize, mimeType, mediaType)
	if err != nil {
		os.Remove(fullPath)
		sendErrorResponse(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}

	mediaID, _ := result.LastInsertId()

	fmt.Printf("💾 [CHUNKED] Сохранено в БД: ID=%d\n", mediaID)

	// Optimize video ASYNCHRONOUSLY if needed
	if mediaType == "video" {
		go func() {
			fmt.Printf("🎬 [ASYNC] Начинается фоновая оптимизация для ID=%d\n", mediaID)

			optimizedPath, err := optimizeVideo(fullPath)
			if err != nil {
				fmt.Printf("❌ [ASYNC] Ошибка оптимизации ID=%d: %v\n", mediaID, err)
				return
			}

			// Update database with optimized file
			optimizedRelativePath := strings.TrimPrefix(optimizedPath, UploadDir+string(filepath.Separator))
			optimizedFileName := filepath.Base(optimizedPath)
			fileInfo, _ := os.Stat(optimizedPath)
			optimizedSize := fileInfo.Size()

			updateQuery := `
				UPDATE user_media 
				SET file_name = ?, file_path = ?, file_size = ?
				WHERE id = ?
			`
			_, err = h.DB.Exec(updateQuery, optimizedFileName, optimizedRelativePath, optimizedSize, mediaID)
			if err != nil {
				fmt.Printf("❌ [ASYNC] Ошибка обновления БД для ID=%d: %v\n", mediaID, err)
				return
			}

			fmt.Printf("✅ [ASYNC] Оптимизация завершена для ID=%d, новый размер: %d\n", mediaID, optimizedSize)
		}()
	}

	// Cleanup temp directory
	os.RemoveAll(uploadDir)

	fmt.Printf("🎉 [CHUNKED] Загрузка завершена: ID=%d (оптимизация в фоне)\n", mediaID)

	sendSuccessResponse(w, map[string]interface{}{
		"id":            int(mediaID),
		"url":           "/api/media/file/" + strconv.Itoa(int(mediaID)),
		"file_name":     finalFileName,
		"original_name": fileName,
		"file_size":     totalSize,
		"mime_type":     mimeType,
		"media_type":    mediaType,
		"optimizing":    mediaType == "video", // Indicate that optimization is in progress
	})
}
