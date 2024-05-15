package repository

import (
	d "bam/internal/core/domain"
	u "bam/internal/core/utils"
	"io"
	"net/http"
	"os"

	"errors"
	"strings"
	"time"

	// "github.com/gofiber/fiber/v2/logger"
	"gorm.io/gorm"
)

var (
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")

	ErrInternal = errors.New("internal server error")
)

type OrbianRepository struct {
	db *gorm.DB
}

func NewOrbianRepository(db *gorm.DB) *OrbianRepository {
	return &OrbianRepository{db}
}

func (r *OrbianRepository) CreateOrbian(orb *d.FormOrbianReq, req *http.Request) (*d.FormOrbianReq, error) {
	// Parse the birth date to extract day, month, and year
	birthday, err := time.Parse("02/01/2006", orb.BirthDay)
	if err != nil {
		return nil, err
	}

	day := birthday.Day()
	month := int(birthday.Month())
	year := birthday.Year()

	// Fetch the lunar date
	monthName := u.GetThaiMonthName(month)
	lunarDate := u.FetchLunarDateAndNaksatNotFormat(day, monthName, year)

	// Update the FormOrbianReq object with the fetched lunar date
	orb.LunarDate = lunarDate

	// Save the image file to the specified path
	file, fileHeader, err := req.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the filename from the multipart.FileHeader
	filename := fileHeader.Filename

	// Create a file in the specified path
	imagePath := "bam/internal/core/upload/" + filename
	outFile, err := os.Create(imagePath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	// Copy the uploaded file to the created file
	_, err = io.Copy(outFile, file)
	if err != nil {
		return nil, err
	}

	// Update the FormOrbianReq object with the image file path
	orb.ImageFilePath = imagePath

	// Create the FormOrbianReq object in the database
	if err := r.db.Create(orb).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ErrConflictingData
		}
		return nil, err // Return the actual error for other cases
	}

	return orb, nil
}
