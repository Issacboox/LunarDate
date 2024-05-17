package model

import "gorm.io/gorm"

type Title string

const (
	Mr   Title = "นาย"
	Mrs  Title = "นาง"
	Miss Title = "นางสาว"
)

// Response represents the outgoing JSON response
type (
	FormOrdianReq struct {
		gorm.Model
		Address1             string `json:"address1"  validate:"required"`
		Address2             string `json:"address2" validate:"required"`
		City                 string `json:"city" validate:"required"`
		Zip                  string `json:"zip" validate:"required"`
		CreateFormDate       string `json:"create_form_date"`
		ImageFilePath        string `json:"-"`
		NameTitle            string `json:"name_title" validate:"required,name_title"`
		FirstName            string `json:"first_name" validate:"required"`
		LastName             string `json:"last_name" validate:"required"`
		IdentityID           string `json:"identity_id" validate:"required,min=13,max=13"`
		FatherTitleName      string `json:"father_tname" validate:"required,name_title"`
		FatherFirstName      string `json:"father_fname" validate:"required"`
		FatherLastName       string `json:"father_lname" validate:"required"`
		MatherTitleName      string `json:"mather_tname" validate:"required,name_title"`
		MatherFirstName      string `json:"mather_fname" validate:"required"`
		MatherLastName       string `json:"mather_lname" validate:"required"`
		BirthDay             string `json:"birth_day" validate:"required,date_format"`
		BirthTime            string `json:"birth_time" validate:"required,birth_time"`
		LunarDate            string `json:"lunar_date"`
		Age                  string `json:"age" validate:"required"`
		Height               string `json:"height" validate:"required"`
		Weight               string `json:"weight" validate:"required"`
		CareerName           string `json:"career_name"`
		WorkingAtCompanyName string `json:"working_at" validate:"required"`
		CompanyPosition      string `json:"company_position" validate:"required"`
		AmountOfOrdians      string `json:"ordian_amount" validate:"required"`
	}

	OrdianResponse struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name" `
		LastName  string `json:"last_name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	GetOrdianResponse struct {
		ID              string `json:"id"`
		NameTitle       string `json:"name_title" `
		FirstName       string `json:"first_name" `
		LastName        string `json:"last_name"`
		Age             string `json:"age" `
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		BirthDay        string `json:"birth_day" `
		BirthTime       string `json:"birth_time" `
		LunarDate       string `json:"lunar_date"`
		Height          string `json:"height"`
		Weight          string `json:"weight" `
		AmountOfOrdians string `json:"ordian_amount" `
	}
)
