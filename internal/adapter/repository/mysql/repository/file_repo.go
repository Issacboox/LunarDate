package repository

import (
	d "bam/internal/core/domain"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {

	return &FileRepository{db}
}

func (r *FileRepository) PostFile(file d.MultiFile) error {
	return r.db.Create(&file).Error
}
