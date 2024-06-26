package service

import (
	d "bam/internal/core/domain"
	p "bam/internal/core/port"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

var (
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")

	ErrInternal = errors.New("internal server error")
)

type OrdianService struct {
	repo p.OrdianRepository
}

func NewOrdianService(ordianRepository p.OrdianRepository) *OrdianService {
	return &OrdianService{ordianRepository}
}

func (ob *OrdianService) OrdianRegister(orb *d.FormOrdianReq, req *http.Request) (*d.FormOrdianReq, error) {
	ordian, err := ob.repo.CreateOrdian(orb, req)
	if err != nil {
		// Log the error for debugging
		log.Println("Error creating ordian:", err)
		return nil, err // Return the actual error
	}
	return ordian, nil
}

func (ob *OrdianService) ListOrdian() ([]*d.FormOrdianReq, error) {
	ordian, err := ob.repo.GetOrdian()
	if err != nil {
		return nil, ErrInternal
	}

	return ordian, nil

}

func (ob *OrdianService) GetOrdianById(ordianId string) (*d.FormOrdianReq, error) {
	ordian, err := ob.repo.GetOrdianById(ordianId)
	if err != nil {
		return nil, ErrInternal
	}

	return ordian, nil
}

func IndentText(text string, indentLevel int) string {
	indent := strings.Repeat("    ", indentLevel)
	return indent + text
}

func (ob *OrdianService) DownloadOrdianByID(id string) ([]byte, error) {
	ordian, err := ob.repo.GetOrdianById(id)
	if err != nil {
		log.Println("Error fetching ordian by ID:", err)
		return nil, ErrInternal
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	fontPath := "C://Users/Sirin/OneDrive/เอกสาร/go/LunarDate/internal/core/utils/pdf/THSarabunNew.ttf"
	pdf.AddUTF8Font("THSarabun", "", fontPath)

	pdf.SetFont("THSarabun", "", 18)
	text := "หนังสือกราบทูลขอประทานการอุปสมบท"
	pdf.CellFormat(0, 10, text, "", 0, "C", false, 0, "")

	// Download the image from the URL
	imageUrl := ordian.ImageFilePath
	resp, err := http.Get(imageUrl)
	if err != nil {
		log.Println("Error downloading image:", err)
		return nil, ErrInternal
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		log.Println("Error determining file extension:", err)
		return nil, ErrInternal
	}

	ext := extensions[0]
	tmpFile, err := ioutil.TempFile("", "image-*"+ext)
	if err != nil {
		log.Println("Error creating temporary file:", err)
		return nil, ErrInternal
	}
	defer os.Remove(tmpFile.Name()) // clean up

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Println("Error saving image to temporary file:", err)
		return nil, ErrInternal
	}

	imagePath := tmpFile.Name()
	imageWidth := 40.0
	imageHeight := 40.0

	// Calculate aspect ratio
	pdfImageOptions := gofpdf.ImageOptions{
		ImageType: "",
		ReadDpi:   true,
	}

	// Get image dimensions
	info := pdf.RegisterImageOptions(imagePath, pdfImageOptions)
	aspectRatio := info.Width() / info.Height()

	// Calculate new dimensions to maintain aspect ratio
	if imageWidth/imageHeight > aspectRatio {
		imageWidth = imageHeight * aspectRatio
	} else {
		imageHeight = imageWidth / aspectRatio
	}

	// Add the image with the calculated dimensions
	pdf.ImageOptions(imagePath, pdf.GetX()-36, pdf.GetY()-4, imageWidth, imageHeight, false, pdfImageOptions, 0, "")
	pdf.Ln(40)

	// Title
	pdf.CellFormat(0, 10, ordian.Address1, "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 10, ordian.Address2, "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 10, ordian.City, "", 1, "R", false, 0, "")

	// Subject
	pdf.SetFont("THSarabun", "", 18)
	pdf.CellFormat(0, 10, "วันที่ "+ordian.CreateFormDate, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Main Content
	pdf.SetFont("THSarabun", "", 18)
	content0 := `เรื่อง     ขอประทานการอุปสมบท`
	content1 := `กราบทูล สมเด็จพระอริยวงศาคตญาณ สมเด็จพระสังฆราช สกลมหาสังฆปริณายก`

	content2 := "   ด้วยเกล้ากระหม่อม " + ordian.NameTitle + ordian.FirstName + " " + ordian.LastName + " เลขที่บัตรประจำตัวประชาชน " + ordian.IdentityID + " เป็นบุตร " + ordian.FatherTitleName + ordian.FatherFirstName + " " + ordian.FatherLastName + " และ " + ordian.MatherTitleName + ordian.MatherFirstName + " " + ordian.MatherLastName + " เกิดเมื่อ " + ordian.LunarDate + " เวลา " + ordian.BirthTime + " ตรงกับวันที่ " + ordian.BirthDay + " อายุ " + ordian.Age + " ปี ส่วนสูง " + ordian.Height + " เซนติเมตร น้ำหนัก " + ordian.Weight + " กิโลกรัม ปัจจุบันประกอบอาชีพ " + ordian.CareerName + " ที่ " + ordian.WorkingAtCompanyName + " ในตำแหน่ง " + ordian.CompanyPosition + " ประสงค์จะอุปสมบทเป็นพระภิกษุในพระพุทธศาสนา โดยขอประทานพระเมตตาฝ่าพระบาทโปรดทรงเป็นพระอุปัชฌายะ"
	content3 := "      " + "การนี้เกล้ากระหม่อมได้ซักซ้อมอุปสมบทวิธีและรายละเอียดเบื้องต้นในการอุปสมบททั้งรับการอบรมกับ"
	content3_1 := "พระมหาคุณณิช คณะไป โดยประสงค์จะอุปสมบทประมาณ " + ordian.AmountOfOrdians + " วัน จักเป็นวันและเวลาใดสุดแต่จะทรงพระกรุณาโปรด"
	content4 := "                จึงกราบทูลมาเพื่อทรงทราบและโปรดประทานอุปสมบท"

	content5 := `ควรมิควรแล้วแต่จะโปรด`
	content6 := `เกล้ากระหม่อม ______________________
    (` + ordian.NameTitle + ` ` + ordian.FirstName + ` ` + ordian.LastName + `)`

	content7 := `อุปสมบทวัน ______________  ที่ _______  เดือน __________ พ.ศ. _________  เวลา _______ น.`
	content8 := `ฉายาว่า ________________________`

	pdf.MultiCell(0, 10, content0, "", "L", false)
	pdf.MultiCell(0, 10, content1, "", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, content2, "", "J", false)
	pdf.MultiCell(0, 10, IndentText(content3, 2)+content3_1, "", "J", false)
	pdf.MultiCell(0, 10, content4, "", "L", false)
	pdf.Ln(10)
	pdf.MultiCell(0, 10, content5, "", "C", false)
	pdf.MultiCell(0, 10, content6, "", "C", false)
	pdf.MultiCell(0, 10, content7, "", "C", false)
	pdf.MultiCell(0, 10, content8, "", "C", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Println("Error generating PDF:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
