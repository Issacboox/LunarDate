package handler

import (
	d "bam/internal/core/domain"
	"bam/internal/core/port"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

type OrbainHandler struct {
	orbianService port.OrbianService
}

func NewOrbianHandler(orbianService port.OrbianService) *OrbainHandler {
	return &OrbainHandler{orbianService}
}

func (h *OrbainHandler) OrbianRegister(c *fiber.Ctx) error {
	// Parse the form data and files
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Get the form data
	data := form.Value

	// Create a FormOrbianReq object from the form data
	orbian := &d.FormOrbianReq{
		Address1:             data["address1"][0],
		Address2:             data["address2"][0],
		City:                 data["city"][0],
		Zip:                  data["zip"][0],
		CreateFormDate:       data["create_form_date"][0],
		NameTitle:            data["name_title"][0],
		FirstName:            data["first_name"][0],
		LastName:             data["last_name"][0],
		IdentityID:           data["identity_id"][0],
		FatherTitleName:      data["father_tname"][0],
		FatherFirstName:      data["father_fname"][0],
		FatherLastName:       data["father_lname"][0],
		MatherTitleName:      data["mather_tname"][0],
		MatherFirstName:      data["mather_fname"][0],
		MatherLastName:       data["mather_lname"][0],
		BirthDay:             data["birth_day"][0],
		BirthTime:            data["birth_time"][0],
		Age:                  data["age"][0],
		Height:               data["height"][0],
		Weight:               data["weight"][0],
		WorkingAtCompanyName: data["working_at"][0],
		CompanyPosition:      data["company_position"][0],
		AmountOfOrdians:      data["ordian_amount"][0],
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the upload directory if it doesn't exist
	uploadDir := "bam/internal/core/upload/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	// Create a new file in the upload directory
	dst, err := os.Create(filepath.Join(uploadDir, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded file to the new file
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	// Set the image file path in the FormOrbianReq object
	orbian.ImageFilePath = filepath.Join(uploadDir, file.Filename)

	// Call the service method to register the orbian
	_, err = h.orbianService.OrbianRegister(orbian,"sssaa")
	if err != nil {
		return err
	}

	// Return a success response
	res := newOrbianResponse(orbian)
	return c.JSON(res)
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

func newOrbianResponse(orb *d.FormOrbianReq) orbianResponse {
	return orbianResponse{
		ID:        strconv.FormatUint(uint64(orb.ID), 10),
		FirstName: orb.FirstName,
		LastName:  orb.LastName,
		CreatedAt: orb.CreatedAt.Format(time.RFC3339), // Convert time.Time to string
		UpdatedAt: orb.UpdatedAt.Format(time.RFC3339), // Convert time.Time to string
	}
}

type orbianResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
