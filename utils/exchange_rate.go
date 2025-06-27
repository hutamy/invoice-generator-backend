package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hutamy/invoice-generator-backend/config"
)

type ExchangeRateResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%s/latest?from=%s&to=%s", cfg.ExchangeRateAPIUrl, fromCurrency, toCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var res ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, err
	}

	rates := res.Rates
	rate, ok := rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s", toCurrency)
	}

	return rate, nil
}
