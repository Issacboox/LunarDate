package handler

import (
	d "bam/internal/core/domain"
	u "bam/internal/core/utils"

	"bam/internal/core/port"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrbianHandler struct {
	orbianService port.OrbianService
}

func NewOrbianHandler(orbianService port.OrbianService) *OrbianHandler {
	return &OrbianHandler{orbianService}
}

func (h *OrbianHandler) OrbianRegister(c *fiber.Ctx) error {
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
		Height:               data["height"][0],
		Weight:               data["weight"][0],
		WorkingAtCompanyName: data["working_at"][0],
		CompanyPosition:      data["company_position"][0],
		AmountOfOrdians:      data["ordian_amount"][0],
	}

	// Use uploadImage function to handle image upload
	imageUrl, err := uploadImage(c)
	if err != nil {
		return err
	}

	// Set the image file path in the FormOrbianReq object
	orbian.ImageFilePath = imageUrl

	// Convert fasthttp.Request to http.Request
	httpReq, err := u.ConvertRequest(c.Request())
	if err != nil {
		return err
	}

	// Call the service method to register the orbian
	_, err = h.orbianService.OrbianRegister(orbian, httpReq)
	if err != nil {
		return err
	}

	// Return a success response
	res := newOrbianResponse(orbian)
	return c.JSON(res)
}

func uploadImage(c *fiber.Ctx) (string, error) {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error in uploading Image:", err)
		return "", fmt.Errorf("server error: %v", err)
	}

	// Check the file extension
	fileExt := filepath.Ext(file.Filename)
	if fileExt != ".jpg" && fileExt != ".png" {
		return "", fmt.Errorf("unsupported file type: %s, only .jpg and .png are allowed", fileExt)
	}

	// Check the file size
	if file.Size > 5*1024*1024 { // 5MB limit
		return "", fmt.Errorf("file size exceeds the limit of 5MB")
	}

	// Generate a unique filename
	uniqueID := uuid.New()
	filename := strings.Replace(uniqueID.String(), "-", "", -1)
	image := fmt.Sprintf("%s%s", filename, fileExt)

	// Set the upload directory
	uploadDir := "C:/Users/Sirin/OneDrive/เอกสาร/go/LunarDate/internal/adapter/repository/upload"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("Error creating upload directory:", err)
		return "", fmt.Errorf("server error: %v", err)
	}
	// Save the file to the upload directory
	imagePath := filepath.Join(uploadDir, image)
	err = c.SaveFile(file, imagePath)
	if err != nil {
		log.Println("Error in saving Image:", err)
		return "", fmt.Errorf("server error: %v", err)
	}

	// Return the URL for accessing the uploaded image
	imageURL := fmt.Sprintf("http://localhost:3000/upload/%s", image)
	return imageURL, nil
}

func newOrbianResponse(orb *d.FormOrbianReq) d.OrbianResponse {
	return d.OrbianResponse{
		ID:        strconv.FormatUint(uint64(orb.ID), 10),
		FirstName: orb.FirstName,
		LastName:  orb.LastName,
		CreatedAt: orb.CreatedAt.Format(time.RFC3339), // Convert time.Time to string
		UpdatedAt: orb.UpdatedAt.Format(time.RFC3339), // Convert time.Time to string
	}
}

func handleError(c *fiber.Ctx, err error) error {
	// Implement your error handling logic here
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}

func handleSuccess(c *fiber.Ctx, res interface{}) error {
	// Implement your success handling logic here
	return c.JSON(res)
}

func (h *OrbianHandler) ListOrbian(c *fiber.Ctx) error {
	orbian, err := h.orbianService.ListOrbian()
	if err != nil {
		return handleError(c, err)
	}

	// Assuming orbianResponse is a struct that can be constructed from FormOrbianReq
	res := make([]d.FormOrbianReq, 0)
	for _, orb := range orbian {
		res = append(res, d.FormOrbianReq{
			// ID:                   orb.ID,
			Address1:             orb.Address1,
			Address2:             orb.Address2,
			City:                 orb.City,
			Zip:                  orb.Zip,
			CreateFormDate:       orb.CreateFormDate,
			NameTitle:            orb.NameTitle,
			FirstName:            orb.FirstName,
			LastName:             orb.LastName,
			IdentityID:           orb.IdentityID,
			FatherTitleName:      orb.FatherTitleName,
			FatherFirstName:      orb.FatherFirstName,
			FatherLastName:       orb.FatherLastName,
			MatherTitleName:      orb.MatherTitleName,
			MatherFirstName:      orb.MatherFirstName,
			MatherLastName:       orb.MatherLastName,
			BirthDay:             orb.BirthDay,
			BirthTime:            orb.BirthTime,
			LunarDate:            orb.LunarDate,
			Age:                  orb.Age,
			Height:               orb.Height,
			Weight:               orb.Weight,
			WorkingAtCompanyName: orb.WorkingAtCompanyName,
			CompanyPosition:      orb.CompanyPosition,
			AmountOfOrdians:      orb.AmountOfOrdians,
		})
	}

	return handleSuccess(c, res)
}
