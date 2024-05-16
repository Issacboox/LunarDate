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

type OrbianRepository struct {
	db *gorm.DB
}

func NewOrbianRepository(db *gorm.DB) *OrbianRepository {
	return &OrbianRepository{db}
}

func (r *OrbianRepository) CreateOrbian(orb *d.FormOrbianReq, req *http.Request) (*d.FormOrbianReq, error) {
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

	// Update the FormOrbianReq object with the fetched lunar date
	orb.LunarDate = lunarDate

	// Create the FormOrbianReq object in the database
	if err := r.db.Create(orb).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ErrConflictingData
		}
		return nil, err // Return the actual error for other cases
	}

	return orb, nil
}

func (r *OrbianRepository) GetOrbian() ([]*d.FormOrbianReq, error) {
    var orbians []*d.FormOrbianReq
    if err := r.db.Find(&orbians).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, ErrDataNotFound
        }
        return nil, err
    }

    return orbians, nil
}