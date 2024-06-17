package service

import (
	d "bam/internal/core/domain"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type LunarDateService interface {
	GetLunarDate(birthday string) (*d.LunarDateResponse, error)
}

type lunarDateServiceImpl struct{}

func NewLunarDateService() LunarDateService {
	return &lunarDateServiceImpl{}
}

const thaiYearOffset = 543

var ThaiMonths = map[int]string{
	1:  "มกราคม",
	2:  "กุมภาพันธ์",
	3:  "มีนาคม",
	4:  "เมษายน",
	5:  "พฤษภาคม",
	6:  "มิถุนายน",
	7:  "กรกฎาคม",
	8:  "สิงหาคม",
	9:  "กันยายน",
	10: "ตุลาคม",
	11: "พฤศจิกายน",
	12: "ธันวาคม",
}

func CalculateLunarDate(date time.Time) (*d.LunarDateResponse, error) {
	thaiMonth := ThaiMonths[int(date.Month())]
	url := fmt.Sprintf("https://www.myhora.com/ปฏิทิน/%d-%s-พ.ศ.%d.aspx", date.Day(), thaiMonth, date.Year()+thaiYearOffset)

	c := colly.NewCollector()

	// Variables to store result
	var lunarDate string

	// Find the Lunar date element
	c.OnHTML("font.f146x", func(e *colly.HTMLElement) {
		lunarDate = e.Text
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error visiting page: %s\n", err)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error visiting URL: %s\n", err)
	}

	if lunarDate == "" {
		return nil, fmt.Errorf("lunar date not found on the page")
	}

	// Remove the prefix "ตรงกับ"
	lunarDate = strings.TrimPrefix(lunarDate, "ตรงกับ")

	// Split the lunar date into three parts
	parts := strings.Split(lunarDate, " ")

	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid lunar date format")
	}

	day := parts[0]
	lunar := strings.Join(parts[1:len(parts)-2], " ")
	naksat := parts[len(parts)-2]
	birthday := date.Format("02-01-2006")

	return &d.LunarDateResponse{
		Birthday:   birthday,
		Age:        calculateAge(date),
		Day:        day,
		LunarDate:  lunar,
		NaksatYear: naksat,
	}, nil
}

// var ThaiLunarMonths = []string{"อ้าย", "ยี่", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า", "สิบ", "สิบเอ็ด", "สิบสอง"}

func ParseThaiDate(dateStr string) (time.Time, error) {
	layout := "02-01-2006"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func FormatThaiDate(date time.Time) string {
	year := date.Year() + thaiYearOffset
	return fmt.Sprintf("%02d/%02d/%04d", date.Day(), date.Month(), year)
}

func calculateAge(birthday time.Time) int {
	today := time.Now()
	age := today.Year() - birthday.Year()
	if today.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

func (s *lunarDateServiceImpl) GetLunarDate(birthday string) (*d.LunarDateResponse, error) {
	date, err := ParseThaiDate(birthday)
	if err != nil {
		return nil, err
	}

	lunarDate, err := CalculateLunarDate(date)
	if err != nil {
		return nil, err
	}

	return lunarDate, nil
}
