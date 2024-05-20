package handler

import (
	d "bam/internal/core/domain"
	u "bam/internal/core/utils"

	"bam/internal/core/port"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrdianHandler struct {
	ordianService port.OrdianService
}

func NewOrdianHandler(ordianService port.OrdianService) *OrdianHandler {
	return &OrdianHandler{ordianService}
}

func getString(data map[string][]string, key string) string {
	if val, ok := data[key]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}

func (h *OrdianHandler) OrdianRegister(c *fiber.Ctx) error {
	// Parse the form data and files
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// Get the form data
	data := form.Value

	ordian := &d.FormOrdianReq{
		Address1:             getString(data, "address1"),
		Address2:             getString(data, "address2"),
		City:                 getString(data, "city"),
		Zip:                  getString(data, "zip"),
		NameTitle:            getString(data, "name_title"),
		FirstName:            getString(data, "first_name"),
		LastName:             getString(data, "last_name"),
		IdentityID:           getString(data, "identity_id"),
		FatherTitleName:      getString(data, "father_tname"),
		FatherFirstName:      getString(data, "father_fname"),
		FatherLastName:       getString(data, "father_lname"),
		MatherTitleName:      getString(data, "mather_tname"),
		MatherFirstName:      getString(data, "mather_fname"),
		MatherLastName:       getString(data, "mather_lname"),
		BirthDay:             getString(data, "birth_day"),
		BirthTime:            getString(data, "birth_time"),
		Height:               getString(data, "height"),
		Weight:               getString(data, "weight"),
		CareerName:           getString(data, "career_name"),
		WorkingAtCompanyName: getString(data, "working_at"),
		CompanyPosition:      getString(data, "company_position"),
		AmountOfOrdians:      getString(data, "ordian_amount"),
	}

	// Use uploadImage function to handle image upload
	imageUrl, err := uploadImage(c)
	if err != nil {
		return err
	}

	ordian.ImageFilePath = imageUrl

	// Convert fasthttp.Request to http.Request
	httpReq, err := u.ConvertRequest(c.Request())
	if err != nil {
		return err
	}

	_, err = h.ordianService.OrdianRegister(ordian, httpReq)
	if err != nil {
		return err
	}

	// Return a success response
	res := newOrdianResponse(ordian)
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
	imageURL := fmt.Sprintf("http://localhost:3000/api/v1/img/%s", image)
	return imageURL, nil
}

func newOrdianResponse(orb *d.FormOrdianReq) d.OrdianResponse {
	return d.OrdianResponse{
		FirstName: orb.FirstName,
		LastName:  orb.LastName,
		CreatedAt: orb.CreatedAt.Format(time.RFC3339),
		UpdatedAt: orb.UpdatedAt.Format(time.RFC3339),
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

func (h *OrdianHandler) ListOrdian(c *fiber.Ctx) error {
	ordian, err := h.ordianService.ListOrdian()
	if err != nil {
		return handleError(c, err)
	}

	res := make([]d.GetOrdianResponse, 0)
	for _, orb := range ordian {
		res = append(res, d.GetOrdianResponse{
			NameTitle:       orb.NameTitle,
			FirstName:       orb.FirstName,
			LastName:        orb.LastName,
			BirthDay:        orb.BirthDay,
			BirthTime:       orb.BirthTime,
			LunarDate:       orb.LunarDate,
			Age:             orb.Age,
			Height:          orb.Height,
			Weight:          orb.Weight,
			AmountOfOrdians: orb.AmountOfOrdians,
		})
	}

	return handleSuccess(c, res)
}

func (h *OrdianHandler) ListOrdianAllData(c *fiber.Ctx) error {
	ordian, err := h.ordianService.ListOrdian()
	if err != nil {
		return handleError(c, err)
	}

	res := make([]d.FormOrdianReq, 0)
	for _, orb := range ordian {
		res = append(res, d.FormOrdianReq{
			Address1:             orb.Address1,
			Address2:             orb.Address2,
			City:                 orb.City,
			Zip:                  orb.Zip,
			CreateFormDate:       orb.CreateFormDate,
			NameTitle:            orb.NameTitle,
			FirstName:            orb.FirstName,
			LastName:             orb.LastName,
			ImageFilePath:        orb.ImageFilePath,
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
			CareerName:           orb.CareerName,
			WorkingAtCompanyName: orb.WorkingAtCompanyName,
			CompanyPosition:      orb.CompanyPosition,
			AmountOfOrdians:      orb.AmountOfOrdians,
		})
	}

	return handleSuccess(c, res)
}

func (h *OrdianHandler) OrdianIdEndpoint(c *fiber.Ctx) error {
	ordianId := c.Params("id")
	if ordianId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request: id is missing")
	}

	ordian, err := h.ordianService.GetOrdianById(ordianId)
	if err != nil {
		return handleError(c, err)
	}

	return handleSuccess(c, ordian)
}

func (h *OrdianHandler) DownloadOrdianByID(c *fiber.Ctx) error {
	id := c.Params("id")

	pdfBytes, err := h.ordianService.DownloadOrdianByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate PDF")
	}

	c.Set(fiber.HeaderContentType, "application/pdf")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename=ordination.pdf")
	return c.Send(pdfBytes)
}
