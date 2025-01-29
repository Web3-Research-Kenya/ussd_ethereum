package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CoinGeckoAPIResponse struct {
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
}

type ExchangeAPIResponse struct {
	Result             string `json:"result"`
	Documentation      string `json:"documentation"`
	TermsOfUse         string `json:"terms_of_use"`
	TimeLastUpdateUnix int    `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string `json:"time_last_update_utc"`
	TimeNextUpdateUnix int    `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string `json:"time_next_update_utc"`
	BaseCode           string `json:"base_code"`
	ConversionRates    struct {
		Usd int     `json:"USD"`
		Jpy float64 `json:"JPY"`
		Kes float64 `json:"KES"`
	} `json:"conversion_rates"`
}

const (
	exchangeRateAPIKey = "d6b7180313f89c7ac86ea784"
	exchangeRateURL    = "https://v6.exchangerate-api.com/v6/%s/latest/%s"

	coingeckoURL    = "https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s"
	coinGeckoAPIKey = "CG-ccGkH1VyyS5sYUbs2QiN2EV1"
)

func ConvertKESToUSD(kes string) (float64, error) {
	kesFloat, err := strconv.ParseFloat(kes, 64)

	if err != nil {
		slog.Error("conversion error", "err", err)
		return 0, err
	}
	url := fmt.Sprintf(exchangeRateURL, exchangeRateAPIKey, "USD")
	agent := fiber.Get(url)

	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		slog.Error("conversion error", "err", errs)
	}

	exchange := new(ExchangeAPIResponse)
	if err := json.Unmarshal(body, exchange); err != nil {
		slog.Error("conversion error", "err", err)
		return 0, err
	}

	converted := int(kesFloat) / exchange.ConversionRates.Usd
	return float64(converted), nil
}

func ConvertUSDToEth(usd float64) (float64, error) {
	url := fmt.Sprintf(coingeckoURL, "ethereum", "usd")
	agent := fiber.Get(url)
	agent.Add("x-cg-demo-api-key", "CG-ccGkH1VyyS5sYUbs2QiN2EV1")

	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		slog.Error("conversion error", "err", errs)
	}

	rates := new(CoinGeckoAPIResponse)
	if err := json.Unmarshal(body, rates); err != nil {
		slog.Error("conversion error", "err", err)
		return 0, err
	}

	converted := usd / rates.Ethereum.Usd
	return converted, nil

}
