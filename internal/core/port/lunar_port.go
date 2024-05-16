package port

import (
	d "bam/internal/core/domain"
)

type LunarRepository interface {
	CreateLunar(lunar *d.FormOrdianReq) (*d.FormOrdianReq, error)
}

type LunarService interface {
	CreateLunar(lunar *d.FormOrdianReq) (*d.FormOrdianReq, error)
}
