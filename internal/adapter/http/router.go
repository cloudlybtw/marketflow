package http

import (
	"net/http"
)

func NewRouter(handler *HTTPHandler) http.Handler {
	mux := http.NewServeMux()
	// mux.HandleFunc("/prices/latest/", handler.GetPriceInfo)
	// mux.HandleFunc("/prices/average/", handler.GetPriceInfo)
	// mux.HandleFunc("/prices/highest/", handler.GetPriceInfo)
	// mux.HandleFunc("/prices/lowest/", handler.GetPriceInfo)
	mux.HandleFunc("POST /mode/test", handler.SwitchMode)
	mux.HandleFunc("POST /mode/live", handler.SwitchMode)
	mux.HandleFunc("/health", handler.Health)
	return mux
}
