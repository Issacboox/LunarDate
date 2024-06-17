package port

import (
	d "bam/internal/core/domain"
)

type LunarDateService interface {
	GetLunarDate(birthday string) (*d.LunarDateResponse, error)
}
