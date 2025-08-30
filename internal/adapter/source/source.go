package source

import (
	"context"
	"log/slog"
	"marketflow/internal/adapter/config"
	"marketflow/internal/core/domain"
	"marketflow/internal/core/service"
	"sync"
)

var Live = true

type ExchangeAdapter interface {
	ReadData(ctx context.Context, exchange config.Exchange, out chan<- domain.Message) error
}

type ModeAdapter struct {
	Live       bool
	exchanges  []config.Exchange
	serv       *service.PriceService
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewModeAdapter(exchanges []config.Exchange, serv *service.PriceService) *ModeAdapter {
	return &ModeAdapter{
		Live:      true,
		exchanges: exchanges,
		serv:      serv,
	}
}

func NewExchangeAdapter(ctx context.Context, exchange config.Exchange, out chan<- domain.Message) ExchangeAdapter {
	if !Live {
		return NewMockExchange([]config.Exchange{exchange}, slog.Default())
	}
	return NewTCPExchange(ctx, exchange, out)
}

func (m *ModeAdapter) Start(ctx context.Context) {
	m.ctx = ctx // Save the context so we can cancel it later
	rawMessages := make(chan domain.Message, 100)
	processedMessages := make(chan domain.Message, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.serv.SavePrice(ctx, processedMessages)
	}()

	for _, exchange := range m.exchanges {
		wg.Add(1)
		go func(ex config.Exchange) {
			defer wg.Done()
			t := NewExchangeAdapter(ctx, ex, rawMessages)
			t.ReadData(ctx, ex, rawMessages)
		}(exchange)
	}

	workerCount := 5
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			service.Run(ctx, workerID, rawMessages, processedMessages)
		}(i)
	}

	<-ctx.Done()
	slog.Info("Context cancelled. Waiting for workers to shut down...")
	wg.Wait()
	slog.Info("All workers shut down")
}

func (m *ModeAdapter) SwitchMode(mode string) {
	var newMode bool
	if mode == "live" {
		newMode = true
	} else {
		newMode = false
	}

	if m.Live == newMode {
		slog.Info("Already in mode:", "mode", mode)
		return
	}

	m.Live = newMode
	Live = newMode

	// Cancel existing pipeline if running
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	// Start new pipeline
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel
	go m.Start(ctx)

	slog.Info("Switched mode", "mode", mode)
}
