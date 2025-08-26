package service

import (
	"errors"
	"strings"
)

var (
	validPairs     = map[string]struct{}{"BTCUSDT": {}, "DOGEUSDT": {}, "TONUSDT": {}, "SOLUSDT": {}, "ETHUSDT": {}}
	validExchanges = map[string]struct{}{"exchange1": {}, "exchange2": {}, "exchange3": {}, "test_exchange1": {}, "test_exchange2": {}, "test_exchange3": {}}
)

func ValidateRequest(exchange, pair string) error {
	if pair != "" {
		if _, ok := validPairs[strings.ToUpper(pair)]; !ok {
			return errors.New("invalid trading pair")
		}
	}

	if exchange != "" {
		if _, ok := validExchanges[strings.ToLower(exchange)]; !ok {
			return errors.New("invalid exchange")
		}
	}

	return nil
}
