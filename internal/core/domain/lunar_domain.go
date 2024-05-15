package model

import "gorm.io/gorm"

type Request struct {
	Birthday string `json:"birthday"`
}

// Response represents the outgoing JSON response
type LunarDate struct {
	gorm.Model
	BirthdayDate string `json:"birthday_date"`
	Day          string `json:"day"`
	LunarDate    string `json:"lunar_date"`
	Naksat       string `json:"naksat"`
}

type CheckLunarDate struct {
	Date      string `json:"date"`
	LunarDate string `json:"lunar_date"`
}
