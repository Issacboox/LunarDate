package port

import (
	d "bam/internal/core/domain"
)

type LunarRepository interface {
	CreateLunar(lunar *d.FormOrbianReq) (*d.FormOrbianReq, error)

}

type LunarService interface {
	CreateLunar(lunar *d.FormOrbianReq) (*d.FormOrbianReq, error)
}