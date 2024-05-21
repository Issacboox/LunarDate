package model

import "gorm.io/gorm"

type FileType string

const (
	MedicalCertificate      FileType = "ใบรับรองแพทย์"
	CrimeCertificate        FileType = "ใบรับรองอาชญากรรม"
	CopyOfIDCard            FileType = "สำเนาบัตรประชาชน"
	CopyOfHouseRegistration FileType = "สำเนาทะเบียนบ้าน"
)

type (
	MultiFile struct {
		gorm.Model
		FileName string   `json:"file_name"`
		FileType FileType `json:"file_type"`
		UserID   int      `json:"user_id"`
		FilePath string   `json:"file_path"`
	}
)
