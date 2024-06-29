package telegramBot

import (
    "log"
    "github.com/Bonittas/ChatBot/openapi"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    if update.Message.IsCommand() {
        switch update.Message.Command() {
        case "start":
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi! I'm a bot that provides exchange rates. Type /rates to see the rates.")
            bot.Send(msg)
        case "rates":
            ratesMsg, err := openapi.GetExchangeRates()
            if err != nil {
                log.Printf("Error fetching exchange rates: %v", err)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't fetch exchange rates at the moment.")
                bot.Send(msg)
                return
            }
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, ratesMsg)
            bot.Send(msg)
        default:
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
            bot.Send(msg)
        }
    } else if strings.Contains(strings.ToLower(update.Message.Text), "exchange rates") {
        ratesMsg, err := openapi.GetExchangeRates()
        if err != nil {
            log.Printf("Error fetching exchange rates: %v", err)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I couldn't fetch exchange rates at the moment.")
            bot.Send(msg)
            return
        }
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, ratesMsg)
        bot.Send(msg)
    } else {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I don't understand that message.")
        bot.Send(msg)
    }
}