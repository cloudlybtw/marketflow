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
	var ex string
	var err error
	switch {
	case exchange == "":
		num, ex, err = s.repo.GetLatest(ctx, pair)
		if err != nil {
			return nil, err
		}
	case exchange != "":
		num, ex, err = s.repo.GetLatestByExchange(ctx, exchange, pair)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request")
	}
	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Exchange:  ex,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (s *PriceService) GetLowest(exchange, pair string, duration time.Duration) (*domain.PriceUpdate, error) {
	var num float64
	var ex string
	var ts int64
	var err error
	switch {
	case exchange == "":
		num, ex, ts, err = s.repo.GetLowestGlobal(pair, duration)
		if err != nil {
			return nil, err
		}
	case exchange != "":
		num, ex, ts, err = s.repo.GetLowestByExchange(exchange, pair, duration)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request")
	}

	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Exchange:  ex,
		Timestamp: ts,
	}, nil
}

func (s *PriceService) GetHighest(exchange, pair string, duration time.Duration) (*domain.PriceUpdate, error) {
	var num float64
	var ex string
	var ts int64
	var err error

	if exchange == "" {
		num, ex, ts, err = s.repo.GetHighestGlobal(pair, duration)
		if err != nil {
			return nil, err
		}
	} else {
		num, ex, ts, err = s.repo.GetHighestByExchange(exchange, pair, duration)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Exchange:  ex,
		Timestamp: ts,
	}, nil
}

func (s *PriceService) GetAverage(exchange, pair string, duration time.Duration) (*domain.PriceUpdate, error) {
	var num float64
	var ex string
	var ts int64
	var err error

	if exchange == "" {
		num, ex, ts, err = s.repo.GetAverageGlobal(pair, duration)
		if err != nil {
			return nil, err
		}
	} else {
		num, ex, ts, err = s.repo.GetAverageByExchange(exchange, pair, duration)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PriceUpdate{
		Symbol:    pair,
		Price:     num,
		Exchange:  ex,
		Timestamp: ts,
	}, nil
}
