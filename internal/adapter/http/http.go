package http

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"marketflow/internal/adapter/source"
	"marketflow/internal/core/domain"
	"marketflow/internal/core/port"
	"marketflow/internal/core/service"
)

type HTTPHandler struct {
	priceService port.PriceService
	mode         *source.ModeAdapter
}

func NewHTTPHandler(priceService port.PriceService, mode *source.ModeAdapter) *HTTPHandler {
	return &HTTPHandler{
		priceService: priceService,
		mode:         mode,
	}
}

func (h *HTTPHandler) GetPriceInfo(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	var info, pair, exchange string
	info = split[2]
	if len(split) == 4 {
		pair, exchange = split[3], ""
	}
	if len(split) == 5 {
		pair, exchange = split[4], split[3]
	}
	if err := service.ValidateRequest(exchange, pair); err != nil {
		Error(w, http.StatusBadRequest, "wrong exchange or pair name")
		return
	}
	var err error
	var duration time.Duration
	params := r.URL.Query()
	period := params.Get("period")
	if period != "" {
		duration, err = time.ParseDuration(period)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if duration > time.Minute*5 {
			Error(w, http.StatusBadRequest, "duration cannot be more than 5 minutes")
			return
		}
	}
	var price *domain.PriceUpdate
	switch info {
	case "latest":
		if duration != 0 {
			Error(w, http.StatusBadRequest, "duration is not applicable for latest price info")
			return
		}
		price, err = h.priceService.GetLatest(r.Context(), exchange, pair)
	// case "average":
	// 	price, err = h.priceService.GetAverage(exchange, pair, duration)
	// case "highest":
	// 	price, err = h.priceService.GetHighest(exchange, pair, duration)
	// case "lowest":
	// 	price, err = h.priceService.GetLowest(exchange, pair, duration)
	default:
		err = errors.New("invalid price info requested")
	}
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, r, price, "GET price info", http.StatusOK)
}

func (h *HTTPHandler) SwitchMode(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	newMode := split[2]
	if (newMode == "live" && h.mode.Live) || (newMode == "test" && !h.mode.Live) {
		Error(w, http.StatusBadRequest, "Mode cannot be changed to current mode")
		return
	}
	// go h.mode.SwitchMode()
	SendResponse(w, r, nil, "Mode changed", http.StatusOK)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	err := h.priceService.CheckHealth(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, r, nil, "Everything is healthy))", http.StatusOK)
}
