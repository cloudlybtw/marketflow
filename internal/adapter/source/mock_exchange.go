package source

import (
	"context"
	"encoding/json"
	"log/slog"
	"marketflow/internal/adapter/config"
	"marketflow/internal/core/domain"
	"time"
)

type MockExchange struct {
	gen *Generator
}

func NewMockExchange(exchanges []config.Exchange, logger *slog.Logger) *MockExchange {
	return &MockExchange{
		gen: NewGenerator(exchanges, logger),
	}
}

func (m *MockExchange) ReadData(ctx context.Context, exchange config.Exchange, out chan<- domain.Message) error {
	
	genOut := make(chan domain.PriceUp, 100)

	
	go m.gen.Start(ctx, genOut)

	for {
		select {
		case <-ctx.Done():
			return nil
		case pu := <-genOut:
			raw, _ := json.Marshal(pu) 
			msg := domain.Message{
				Exchange:  exchange.Name,
				Timestamp: time.Now(), 
				Raw:       raw,
			}
			select {
			case out <- msg:
			case <-ctx.Done():
				return nil
			default:
			
			}
		}
	}
}
