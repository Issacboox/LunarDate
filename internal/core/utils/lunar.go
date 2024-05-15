package utils

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func FetchLunarDateAndNaksat(day int, monthName string, year int) (string, string, string) {
	c := colly.NewCollector()

	var lunarDate, dayName, naksat string

	c.OnHTML("font.f146x", func(e *colly.HTMLElement) {
		html := e.DOM.Text()
		text := strings.TrimSpace(strings.ReplaceAll(html, "\n", ""))

		// Remove "ตรงกับ" if present
		text = strings.Replace(text, "ตรงกับ", "", -1)

		// Split the text to extract needed parts
		parts := strings.Split(text, " ")
		if len(parts) >= 7 {
			dayName = parts[0]
			lunarDate = strings.Join(parts[1:len(parts)-1], " ")
			naksat = parts[len(parts)-1]
		}
	})

	url := fmt.Sprintf("https://www.myhora.com/ปฏิทิน/%d-%s-พ.ศ.%d.aspx", day, monthName, year)
	c.Visit(url)

	// Wait for the scraping to finish
	c.Wait()

	return lunarDate, dayName, naksat
}

func GetThaiMonthName(month int) string {
	months := [...]string{
		"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
		"กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม",
	}
	return months[month-1]
}

func FetchLunarDateAndNaksatNotFormat(day int, monthName string, year int) string {
	c := colly.NewCollector()

	var lunarDate string

	c.OnHTML("font.f146x", func(e *colly.HTMLElement) {
		html := e.DOM.Text()
		lunarDate = strings.TrimSpace(html)
		// Remove the text "ตรงกับวัน" from the lunar date
		lunarDate = strings.Replace(lunarDate, "ตรงกับวัน", "", -1)
	})

	url := fmt.Sprintf("https://www.myhora.com/ปฏิทิน/%d-%s-พ.ศ.%d.aspx", day, monthName, year)
	c.Visit(url)

	c.Wait()

	return lunarDate
}
