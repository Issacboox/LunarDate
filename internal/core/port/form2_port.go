package port

import (
	d "bam/internal/core/domain"
	"net/http"
)

type OrbianRepository interface {
	CreateOrbian(orb *d.FormOrbianReq, req *http.Request) (*d.FormOrbianReq, error)
}

type OrbianService interface {
	OrbianRegister(orb *d.FormOrbianReq, req *http.Request) (*d.FormOrbianReq, error)
}
