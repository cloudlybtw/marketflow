package source

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"marketflow/internal/adapter/config"
	"marketflow/internal/core/domain"
)

type Generator struct {
	exchanges  []config.Exchange
	logger     *slog.Logger
	basePrices map[string]float64
}

func NewGenerator(exchanges []config.Exchange, logger *slog.Logger) *Generator {
	return &Generator{
		exchanges: exchanges,
		logger:    logger,
		basePrices: map[string]float64{
			"BTCUSDT":  43250.50,
			"ETHUSDT":  2520.30,
			"DOGEUSDT": 0.085,
			"TONUSDT":  2.45,
			"SOLUSDT":  98.75,
		},
	}
}

func (g *Generator) Start(ctx context.Context, output chan<- domain.PriceUpdate) {
	g.logger.Info("Starting test data generator")

	symbols := []string{"BTCUSDT", "ETHUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT"}
	exchanges := []string{"exchange1", "exchange2", "exchange3"}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.logger.Info("Test data generator stopping")
			return
		case <-ticker.C:
			for _, symbol := range symbols {
				for _, exchange := range exchanges {
					update := g.generatePriceUpdate(symbol, exchange)
					select {
					case output <- update:
					case <-ctx.Done():
						return
					default:
						g.logger.Warn("Output channel full, dropping test update")
					}
				}
			}
		}
	}
}

func (g *Generator) generatePriceUpdate(symbol, exchange string) domain.PriceUpdate {
	basePrice := g.basePrices[symbol]

	change := (rand.Float64() - 0.5) * 0.02
	newPrice := basePrice * (1 + change)

	g.basePrices[symbol] = newPrice

	return domain.PriceUpdate{
		Symbol:    symbol,
		Price:     newPrice,
		Exchange:  exchange,
		Timestamp: time.Now().Unix(),
	}
}
