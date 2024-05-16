package utils

import (
	d "bam/internal/core/domain"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func CalculateAge(birthDay string) (int, error) {
	// วันเกิด
	birthday, err := time.Parse("02/01/2006", birthDay)
	if err != nil {
		return 0, err
	}

	// วันปัจจุบัน
	now := time.Now()

	// คำนวณอายุ
	age := (now.Year() + 543) - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}

	return age, nil
}

func GetCurrentThaiDate() string {
	// Get the current date
	currentDate := time.Now()

	// Convert the current date to Thai Buddhist Era (พ.ศ.)
	yearThai := currentDate.Year() + 543

	// Format the date as DD/MM/YYYY in Thai Buddhist Era (พ.ศ.)
	return currentDate.Format("02/01/") + strconv.Itoa(yearThai)
}

func IsTitleNameCorrect(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return title == string(d.Mr) || title == string(d.Mrs) || title == string(d.Miss)
}

func FormatDateCorrect(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	layout := "02/01/2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		return false
	}
	return t.Before(time.Now())
}

func FormatTimeBirth(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	layout := "15:04"
	_, err := time.Parse(layout, timeStr)
	return err == nil
}
