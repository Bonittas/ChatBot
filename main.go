package main

import (
    "log"
    "github.com/Bonittas/ChatBot/telegramBot"
    // "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botToken = "6511406154:AAEeyiHaG_PZxF6jYyp38aG370tbr_BTlqI"

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

        telegramBot.HandleUpdate(bot, update)
    }
}
