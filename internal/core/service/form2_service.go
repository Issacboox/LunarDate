package service

import (
	d "bam/internal/core/domain"
	p "bam/internal/core/port"
	"bytes"
	"errors"
	"log"
	"net/http"

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

func (ob *OrdianService) DownloadOrdianByID(id string) ([]byte, error) {
	ordian, err := ob.repo.GetOrdianById(id)
	if err != nil {
		return nil, ErrInternal
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("THSarabun", "", "THSarabunNew.ttf")
	pdf.SetFont("THSarabun", "", 16)

	// Title
	pdf.CellFormat(0, 10, "วันที่ "+ordian.CreateFormDate, "", 1, "R", false, 0, "")
	pdf.Ln(10)

	// Subject
	pdf.SetFont("THSarabun", "", 14)
	pdf.CellFormat(0, 10, "เรื่อง    ขอประทานการอุปสมบท", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Main Content
	pdf.SetFont("THSarabun", "", 14)
	content := `กราบทูล สมเด็จพระอริยวงศาคตญาณ สมเด็จพระสังฆราช สกลมหาสังฆปริณายก
	ด้วยเกล้ากระหม่อม ` + ordian.NameTitle + ` ` + ordian.FirstName + ` ` + ordian.LastName + ` เลขที่บัตรประจำตัวประชาชน ` + ordian.IdentityID + ` เป็นบุตร ` + ordian.FatherTitleName + ` ` + ordian.FatherFirstName + ` ` + ordian.FatherLastName + ` และ ` + ordian.MatherTitleName + ` ` + ordian.MatherFirstName + ` ` + ordian.MatherLastName + ` เกิดเมื่อ ` + ordian.BirthDay + ` เวลา ` + ordian.BirthTime + ` อายุ ` + ordian.Age + ` ปี ส่วนสูง ` + ordian.Height + ` เซนติเมตร น้ำหนัก ` + ordian.Weight + ` กิโลกรัม ปัจจุบันประกอบอาชีพ ` + ordian.CareerName + ` ที่ ` + ordian.WorkingAtCompanyName + ` ในตำแหน่ง ` + ordian.CompanyPosition + ` ประสงค์จะอุปสมบทเป็นพระภิกษุในพระพุทธศาสนา โดยขอประทานพระเมตตาฝ่าพระบาทโปรดทรงเป็นพระอุปัชฌายะ

	การนี้ เกล้ากระหม่อมได้ซักซ้อมอุปสมบทวิธีและรายละเอียดเบื้องต้นในการอุปสมบท ทั้งรับการอบรมกับพระมหาคุณณิช คณะไป โดยประสงค์จะอุปสมบทประมาณ ` + ordian.AmountOfOrdians + ` วัน จักเป็นวันและเวลาใดสุดแต่จะทรงพระกรุณาโปรด

	จึงกราบทูลมาเพื่อทรงทราบและโปรดประทานอุปสมบท

	ควรมิควรแล้วแต่จะโปรด

	เกล้ากระหม่อม ______________________
	(` + ordian.NameTitle + ` ` + ordian.FirstName + ` ` + ordian.LastName + `)`

	pdf.MultiCell(0, 10, content, "", "L", false)

	// Create a buffer to store the PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	// Return the PDF as a byte slice
	return buf.Bytes(), nil
}
