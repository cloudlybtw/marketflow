package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"marketflow/internal/core/domain"
)

// AggregateAndSaveByExchange
func (s *PriceService) SavePrice(ctx context.Context, in <-chan domain.Message) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	aggregation := domain.AggregatedData{
		Exchanges: make(map[string]*domain.ExchangeData),
	}

	for {
		select {
		case <-ctx.Done():
			if len(aggregation.Exchanges) > 0 {
				aggregation.WindowEnd = time.Now()
				s.repo.SaveAggregated(aggregation)
			}
			return

		case msg := <-in:
			var priceUpdate domain.PriceUpdate
			if err := json.Unmarshal(msg.Raw, &priceUpdate); err != nil {
				slog.Info(fmt.Sprintf("Error unmarshaling price update: %v", err))
				continue
			}
			s.repo.SaveLatest(ctx, msg.Exchange, priceUpdate.Symbol, priceUpdate.Price)

			if !isTargetCoin(priceUpdate.Symbol) {
				continue
			}

			// Initialize time window if first message
			if aggregation.WindowStart.IsZero() {
				aggregation.WindowStart = time.Now()
			}

			// Get or create exchange data
			exchangeData, exists := aggregation.Exchanges[msg.Exchange]
			if !exists {
				exchangeData = &domain.ExchangeData{
					Coins: make(map[string]*domain.CoinStats),
				}
				aggregation.Exchanges[msg.Exchange] = exchangeData
			}

			// Get or create coin stats within this exchange
			coinStats, exists := exchangeData.Coins[priceUpdate.Symbol]
			if !exists {
				coinStats = &domain.CoinStats{
					MinPrice:  priceUpdate.Price,
					MaxPrice:  priceUpdate.Price,
					FirstSeen: time.Now(),
				}
				exchangeData.Coins[priceUpdate.Symbol] = coinStats
			}

			// Update statistics
			if coinStats.Count == 0 || priceUpdate.Price < coinStats.MinPrice {
				coinStats.MinPrice = priceUpdate.Price
			}
			if coinStats.Count == 0 || priceUpdate.Price > coinStats.MaxPrice {
				coinStats.MaxPrice = priceUpdate.Price
			}

			// Calculate running average
			coinStats.AvgPrice = (coinStats.AvgPrice*float64(coinStats.Count) + priceUpdate.Price) / float64(coinStats.Count+1)
			coinStats.LastPrice = priceUpdate.Price
			coinStats.Count++
			coinStats.LastSeen = time.Now()

		case <-ticker.C:
			if len(aggregation.Exchanges) > 0 {
				aggregation.WindowEnd = time.Now()
				s.repo.SaveAggregated(aggregation)

				// Reset for new window
				aggregation = domain.AggregatedData{
					WindowStart: time.Now(),
					Exchanges:   make(map[string]*domain.ExchangeData),
				}
			}
		}
	}
}

func isTargetCoin(symbol string) bool {
	targetCoins := map[string]bool{
		"BTCUSDT":  true,
		"DOGEUSDT": true,
		"TONUSDT":  true,
		"SOLUSDT":  true,
		"ETHUSDT":  true,
	}
	return targetCoins[symbol]
}
