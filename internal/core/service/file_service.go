package service

import (
	d "bam/internal/core/domain"
	p "bam/internal/core/port"

	"github.com/jung-kurt/gofpdf"
)

type FileService struct {
	repo p.FileRepository
}

func NewFileService(fileRepository p.FileRepository) *FileService {
	return &FileService{fileRepository}
}

func (s *FileService) AttachFile(file d.MultiFile) error {
	return s.repo.PostFile(file)
}

func (s *FileService) ConvertImageToPDF(imagePath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.Image(imagePath, 10, 10, 200, 0, false, "", 0, "")

	pdfPath := imagePath + ".pdf"
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", err
	}

	return pdfPath, nil
}
