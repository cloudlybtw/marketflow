package service

import (
	"context"

	"marketflow/internal/core/port"
)

type PriceService struct {
	repo port.PriceRepository
}

func NewPriceService(r port.PriceRepository) *PriceService {
	return &PriceService{
		repo: r,
	}
}

func (s *PriceService) CheckHealth(ctx context.Context) error {
	return s.repo.Health(ctx)
}
