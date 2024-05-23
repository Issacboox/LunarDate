package handler

import (
	"fmt"
	"os"
	"path/filepath"

	d "bam/internal/core/domain"
	p "bam/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileService p.FileService
}

func NewFileHandler(fileService p.FileService) *FileHandler {
	return &FileHandler{fileService}
}

func (h *FileHandler) UploadFiles(c *fiber.Ctx) error {
	userID := 1 // Assuming user ID is 1

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	files := form.File["files"]
	fileTypes := form.Value["file_type"]
	if len(files) != len(fileTypes) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Files and file types count mismatch"})
	}

	typeTracker := make(map[d.FileType]bool)
	for i, file := range files {
		fileType := d.FileType(fileTypes[i])
		if typeTracker[fileType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Duplicate file type: %s", fileType)})
		}
		typeTracker[fileType] = true

		// Create the uploads directory if it doesn't exist
		if err := os.MkdirAll("uploads", 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create uploads directory: %s", err.Error())})
		}

		filePath := filepath.Join("uploads", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to save file: %s", err.Error())})
		}

		var pdfPath string
		if filepath.Ext(file.Filename) != ".pdf" {
			pdfPath, err = h.fileService.ConvertImageToPDF(filePath)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to convert image to PDF: %s", err.Error())})
			}

			// Update the file path to point to the PDF file
			filePath = pdfPath
		}

		// Update the file path to include the base URL
		filePath = "http://127.0.0.1:3000/api/v1/file/" + filepath.Base(filePath)

		newFile := d.MultiFile{
			FileName: file.Filename,
			FileType: fileType,
			UserID:   userID,
			FilePath: filePath,
		}

		if err := h.fileService.AttachFile(newFile); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to save file info to database: %s", err.Error())})
		}
	}

	return c.JSON(fiber.Map{"message": "Files uploaded successfully"})
}
