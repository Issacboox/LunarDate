package model

type (
	Request struct {
		Birthday string `json:"birthday"`
	}
	LunarDateResponse struct {
		Birthday   string `json:"birthday"`
		Age        int    `json:"age"`
		Day        string `json:"day"`
		LunarDate  string `json:"lunar_date"`
		NaksatYear string `json:"naksat_year"`
	}
	CheckLunarDate struct {
		Date      string `json:"date"`
		LunarDate string `json:"lunar_date"`
	}
)
