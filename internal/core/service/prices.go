package service

import (
	"context"
	"errors"
	"time"

	"marketflow/internal/core/domain"
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

func (s *PriceService) GetLatest(ctx context.Context, exchange string, pair string) (*domain.PriceUpdate, error) {
	var num float64
	var err error
	switch {
	case exchange == "":
		num, err = s.repo.GetLatest(ctx, pair)
		if err != nil {
			return nil, err
		}
	case exchange != "":
		num, err = s.repo.GetLatestByExchange(ctx, exchange, pair)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request")
	}
	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (s *PriceService) GetLowest(exchange, pair string, duration time.Duration) (*domain.PriceUpdate, error) {
	var num float64
	var err error
	switch {
	case exchange == "":
		return nil, errors.New("temporary unavailable. get lowest")
	case exchange != "":
		num, err = s.repo.GetLowestByExchange(exchange, pair, duration)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request")
	}

	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Timestamp: time.Now().Unix(),
	}, nil
}
