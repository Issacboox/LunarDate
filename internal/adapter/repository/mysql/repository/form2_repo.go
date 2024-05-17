package repository

import (
	d "bam/internal/core/domain"
	u "bam/internal/core/utils"
	"net/http"
	"strconv"

	"errors"
	"strings"
	"time"

	// "github.com/gofiber/fiber/v2/logger"
	"gorm.io/gorm"
)

var (
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")

	ErrInternal = errors.New("internal server error")

	ErrDataNotFound = errors.New("data not found")
)

type OrdianRepository struct {
	db *gorm.DB
}

func NewOrdianRepository(db *gorm.DB) *OrdianRepository {
	return &OrdianRepository{db}
}

func (r *OrdianRepository) CreateOrdian(orb *d.FormOrdianReq, req *http.Request) (*d.FormOrdianReq, error) {
	// Calculate the age
	age, err := u.CalculateAge(orb.BirthDay)
	if err != nil {
		return nil, err
	}
	orb.Age = strconv.Itoa(age)
	orb.CreateFormDate = u.GetCurrentThaiDate()

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

	orb.LunarDate = lunarDate

	if err := r.db.Create(orb).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ErrConflictingData
		}
		return nil, err // Return the actual error for other cases
	}

	return orb, nil
}

func (r *OrdianRepository) GetOrdian() ([]*d.FormOrdianReq, error) {
	var ordians []*d.FormOrdianReq
	if err := r.db.Find(&ordians).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDataNotFound
		}
		return nil, err
	}

	return ordians, nil
}

func (r *OrdianRepository) GetOrdianById(ordianId string) (*d.FormOrdianReq, error) {
	var ordian *d.FormOrdianReq
	if err := r.db.Where("id = ?", ordianId).Find(&ordian).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	return ordian, nil
}

func (r *OrdianRepository) GetOrdianByIdPDF(ordianId string) (*d.FormOrdianReq, error) {
	var ordian d.FormOrdianReq
	if err := r.db.Where("id = ?", ordianId).Find(&ordian).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	return &ordian, nil
}
