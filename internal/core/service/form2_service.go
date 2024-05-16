package service

import (
	d "bam/internal/core/domain"
	p "bam/internal/core/port"
	"errors"
	"log"
	"net/http"
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

func (ob *OrdianService) GetOrdianById(ordianId string) ([]*d.FormOrdianReq, error) {
	ordian, err := ob.repo.GetOrdianById(ordianId)
	if err != nil {
		return nil, ErrInternal
	}

	return ordian, nil
}
