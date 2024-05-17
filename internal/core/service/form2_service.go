package service

import (
	d "bam/internal/core/domain"
	p "bam/internal/core/port"
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"

	                                "github.com/signintech/gopdf"
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
func renderTemplate(tmpl string, data d.FormOrdianReq) (string, error) {
	t, err := template.New("pdf").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func htmlToPDF(html string) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	// pdf.SetFont("Arial", "", 14)
	pdf.Cell(nil, html)

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (ob *OrdianService) DownloadOrdianByID(id string) ([]byte, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link
      rel="stylesheet"
      href="https://fonts.googleapis.com/css2?family=Angsana+New&display=swap"
    />
    <link rel="stylesheet" href="/style.css" />
    <title>Document</title>
  </head>
  <style>
      body {
      background: rgb(204, 204, 204);
      }
      page {
      background: white;
      display: block;
      margin: 0 auto;
      margin-bottom: 0.5cm;
      box-shadow: 0 0 0.5cm rgba(0, 0, 0, 0.5);
      font-family: "Angsana New", Arial, sans-serif;
      position: absolute;
      left: 0;
      right: 0;
      }

      page[size="A4"] {
      width: 21cm;
      height: 29.7cm;
      
      }
      page[size="A4"][layout="landscape"] {
      width: 29.7cm;
      height: 21cm;
      }
      @media print {
      body,
      page {
        background: white;
        margin: 0;
        box-shadow: 0;
      }
      }
      div.top-form {
      display: flex;
      gap: 70px;
      justify-content: center;
      margin-left: 100px;
      }
      div.top-form > h1.title {
      font-size: 1.9em;
      font-weight: 500;
      text-align: center;
      margin-top: 80px;
      margin-left: 10px;
      margin-bottom: 20px;
      }
      img.img-ordian {
      width: 2.8cm;
      height: 3.5cm;
      margin-top: 50px;
      background-color: white;
      border: 0px solid black;
      object-fit: cover;
      }
      div.address {
      display: flex;
      justify-content: end;
      text-align: end;
      margin-right: 80px;
      }
      div.address > h2 {
      font-weight: 500;
      }
      div.created-form-date > h2{
      font-weight: 500;
      text-align: center;
      margin-left:110px;
      }
      div.created-form-header > h2 {
      font-weight: 500;
      margin-left: 100px;
      margin-right: 100px;
      text-align: justify;
      }
      div.form-sign > h2 {
      font-weight: 500;
      margin-left: 100px;
      margin-right: 100px;
      text-align: center;
      }
      div.ordian > h2 {
      font-weight: 500;
      margin-left: 100px;
      margin-right: 100px;
      text-align: justify;
      }
      div.footer > h2 {
      font-weight: 500;
      text-align: center;
      margin-left: 100px;
      margin-right: 100px;
      }
  </style>
  <body>
  {{range .}}
    <div class="">
      <page size="A4">
        <div class="top-form">
          <h1 class="title">หนังสือกราบทูลขอประทานการอุปสมบท</h1>
          <img
            src="https://img.freepik.com/free-photo/portrait-optimistic-businessman-formalwear_1262-3600.jpg?size=626&ext=jpg&ga=GA1.1.553209589.1714953600&semt=ais"
            class="img-ordian"
          />
        </div>
        <div class="address">
          <h2>
      {{ .Address1 }} <br />
          {{ .Address2 }} <br />
      {{ .City }} {{ .Zip }}
          </h2>
        </div>
        <div class="created-form-date">
          <h2>วันที่ {{ .CreateFormDate }}</h2>
        </div>
        <div class="created-form-header">
          <h2>
            เรื่อง&emsp;&emsp;ขอประทานการอุปสมบท
            <br />กราบทูล&emsp;สมเด็จพระอริยวงศาคตญาณ สมเด็จพระสังฆราช
            สกลมหาสังฆปริณายก
          </h2>
          <h2>
            &emsp;&emsp;ด้วยเกล้ากระหม่อม {{ .NameTitle }}{{ .FirstName }}  {{ .LastName }} เลขที่บัตรประจำตัวประชาชน
            {{ .IdentityID }} เป็นบุตร {{ .FatherTitleName }}{{ .FatherFirstName }}  {{ .FatherLastName }} และ {{ .MatherTitleName }}{{ .MatherFirstName }}  {{ .MatherLastName }}
            เกิดเมื่อวัน{{ .LunarDate }} เวลา {{ .BirthTime }} น.
            ตรงกับวันที่ {{ .BirthDay }} อายุ {{ .Age }} ปี ส่วนสูง {{ .Height }}
            เซนติเมตร น้ำหนัก {{ .Weight }} กิโลกรัม ปัจจุบันประกอบอาชีพ{{ .LunarDate }}
            ที่ชื่อสถานที่ทำงาน {{ .WorkingAtCompanyName }} ในตำแหน่งตำแหน่งงาน {{ .CompanyPosition }}
            ประสงค์จะอุปสมบทเป็นพระภิกษุในพระพุทธศาสนา
            โดยขอประทานพระเมตตาฝ่าพระบาทโปรดทรงเป็นพระอุปัชฌายะ
          </h2>
          <h2>
            &emsp;&emsp;การนี้
            เกล้ากระหม่อมได้ซักซ้อมอุปสมบทวิธีและรายละเอียดเบื้องต้นในการอุปสมบท
            ทั้งรับการอบรมกับพระมหาคณิศร คณะใน โดยประสงค์จะอุปสมบทประมาณ {{ .AmountOfOrdians }} วัน
            จักเป็นวันและเวลาใดสุดแต่จะทรงพระกรุณาโปรด<br />
            &emsp;&emsp;จึงกราบทูลมาเพื่อทรงทราบและประทานอุปสมบท
          </h2>
        </div>
        <div class="form-sign">
          <h2>
            ควรมิควรแล้วแต่จะโปรด<br />
            เกล้ากระหม่อม ________________________________<br />
            ({{ .NameTitle }}{{ .FirstName }}  {{ .LastName }})
          </h2>
        </div>
        <br/>
        <div class="ordian">
          <h2>
            อุปสมบทวัน _______ ที่ ___ เดือน ____________________ พ.ศ.
            _________เวลา ______น.
          </h2>
        </div>
        <div class="footer">
            <h2>
                ฉายาว่า
                ________________________________
        </h2>
        </div>
        </page>
      </div>
      </body>
      {{end}} 
      </html>
    `

    data, err := ob.GetOrdianById(id)
    if err != nil {
        return nil, err
    }

    // Step 2: Render HTML Template
    html, err := renderTemplate(tmpl, *data) // Pass data as a slice
    if err != nil {
        return nil, err
    }

    // Step 3: Convert HTML to PDF
    pdfBytes, err := htmlToPDF(html)
    if err != nil {
        return nil, err
    }

    return pdfBytes, nil
}