package telegram

import (
    "log"
    "strings"
    "github.com/Bonittas/ChatBot/openapi"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(botToken, apiUrl string) error {
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        return err
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            if update.Message.IsCommand() {
                switch update.Message.Command() {
                case "start":
                    welcomeMessage := "Welcome to the Questionary Bot! Ask me any question and I'll fetch the answer for you."
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
                    bot.Send(msg)
                default:
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I don't understand that command.")
                    bot.Send(msg)
                }
            } else {
                response := handleMessage(update.Message.Text, apiUrl)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
                bot.Send(msg)
            }
        }
    }
    return nil
}
func handleMessage(message, apiUrl string) string {
    client := openapi.NewClient(apiUrl)
    response, err := client.GetAnswer(strings.TrimSpace(message))
    if err != nil {
        log.Printf("Error getting answer: %v", err)
        return "Sorry, I encountered an error while trying to get an answer for you."
    }

    if response == "" {
        return "I'm sorry, I don't have an answer for that question."
    }

    return response
}