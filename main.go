package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
    botToken         = "6511406154:AAEeyiHaG_PZxF6jYyp38aG370tbr_BTlqI"
    exchangeRateAPI  = "https://v6.exchangerate-api.com/v6/%s/latest/%s"
    apiKey           = "adf234c8fe244535714b13ef"
    baseCurrency     = "USD" 
)

type ExchangeRateResponse struct {
    ConversionRates map[string]float64 `json:"conversion_rates"`
    Result          string             `json:"result"`
    Documentation   string             `json:"documentation"`
    TermsOfUse      string             `json:"terms_of_use"`
    TimeLastUpdate  string             `json:"time_last_update_utc"`
}

func main() {
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        if update.Message.IsCommand() {
            switch update.Message.Command() {
            case "start":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi! I'm a bot that provides exchange rates. Type /rates to see the rates.")
                bot.Send(msg)
            case "rates":
                ratesMsg, err := getExchangeRates()
                if err != nil {
                    log.Printf("Error fetching exchange rates: %v", err)
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't fetch exchange rates at the moment.")
                    bot.Send(msg)
                    continue
                }
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, ratesMsg)
                bot.Send(msg)
            default:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
                bot.Send(msg)
            }
        } else if strings.Contains(strings.ToLower(update.Message.Text), "exchange rates") {
            ratesMsg, err := getExchangeRates()
            if err != nil {
                log.Printf("Error fetching exchange rates: %v", err)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't fetch exchange rates at the moment.")
                bot.Send(msg)
                continue
            }
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, ratesMsg)
            bot.Send(msg)
        } else {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I don't understand that message.")
            bot.Send(msg)
        }
    }
}

func getExchangeRates() (string, error) {
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
