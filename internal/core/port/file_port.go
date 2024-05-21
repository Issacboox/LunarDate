package port

import (
	d "bam/internal/core/domain"
)

type FileRepository interface {
	PostFile(d.MultiFile) error
}

type FileService interface {
	AttachFile(file d.MultiFile) error
	ConvertImageToPDF(imagePath string) (string, error)
}
