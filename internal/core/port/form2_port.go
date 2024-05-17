package port

import (
	d "bam/internal/core/domain"
	"net/http"
)

type OrdianRepository interface {
	CreateOrdian(orb *d.FormOrdianReq, req *http.Request) (*d.FormOrdianReq, error)
	GetOrdian() ([]*d.FormOrdianReq, error)
	GetOrdianById(ordianId string) (*d.FormOrdianReq, error)

}

type OrdianService interface {
	OrdianRegister(orb *d.FormOrdianReq, req *http.Request) (*d.FormOrdianReq, error)
	ListOrdian() ([]*d.FormOrdianReq, error)
	GetOrdianById(ordianId string) (*d.FormOrdianReq, error)

	DownloadOrdianByID(id string) ([]byte, error)
}
