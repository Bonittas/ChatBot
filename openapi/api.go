package openapi

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

const (
    exchangeRateAPI = "https://v6.exchangerate-api.com/v6/%s/latest/%s"
    apiKey          = "adf234c8fe244535714b13ef"
    baseCurrency    = "USD"
)

type ExchangeRateResponse struct {
    ConversionRates map[string]float64 `json:"conversion_rates"`
    Result          string             `json:"result"`
    Documentation   string             `json:"documentation"`
    TermsOfUse      string             `json:"terms_of_use"`
    TimeLastUpdate  string             `json:"time_last_update_utc"`
}

func GetExchangeRates() (string, error) {
    apiURL := fmt.Sprintf(exchangeRateAPI, apiKey, baseCurrency)

    resp, err := http.Get(apiURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var exchangeRates ExchangeRateResponse
    if err := json.NewDecoder(resp.Body).Decode(&exchangeRates); err != nil {
        return "", err
    }

    if exchangeRates.Result != "success" {
        return "", fmt.Errorf("API request failed: %s", exchangeRates.Result)
    }

    ratesMsg := fmt.Sprintf("Exchange Rates (Base Currency: %s):\n", baseCurrency)
    for currency, rate := range exchangeRates.ConversionRates {
        ratesMsg += fmt.Sprintf("%s: %.2f\n", strings.ToUpper(currency), rate)
    }
    return ratesMsg, nil
}