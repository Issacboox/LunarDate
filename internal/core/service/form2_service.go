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

type OrbianService struct {
	repo p.OrbianRepository
}

func NewOrbianService(orbianRepository p.OrbianRepository) *OrbianService {
	return &OrbianService{orbianRepository}
}

func (ob *OrbianService) OrbianRegister(orb *d.FormOrbianReq, req *http.Request) (*d.FormOrbianReq, error) {
	orbian, err := ob.repo.CreateOrbian(orb, req)
	if err != nil {
		// Log the error for debugging
		log.Println("Error creating orbian:", err)
		return nil, err // Return the actual error
	}
	return orbian, nil
}
