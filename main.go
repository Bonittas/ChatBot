package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	telegramBotToken = "6511406154:AAEeyiHaG_PZxF6jYyp38aG370tbr_BTlqI" // Replace with your Telegram bot token
	openTriviaURL    = "https://opentdb.com/api.php?amount=1&type=boolean" // Example API endpoint
	baseURL          = "https://api.telegram.org/bot%s/"
	getUpdatesURL    = "getUpdates?offset=%d&timeout=60" // Timeout set to 60s
)

// TelegramUpdate represents a Telegram update received from polling
type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text     string `json:"text"`
		Chat     Chat   `json:"chat"`
		From     User   `json:"from"`
		MessageID int    `json:"message_id"`
	} `json:"message"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID int `json:"id"`
}

// User represents a Telegram user
type User struct {
	ID int `json:"id"`
}

// TelegramResponse represents a Telegram API response
type TelegramResponse struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func main() {
	offset := 0
	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
			continue
		}

		for _, update := range updates {
			offset = update.UpdateID + 1

			if strings.Contains(strings.ToLower(update.Message.Text), "/start") {
				sendMessage(update.Message.Chat.ID, "Hello! I'm a simple bot. Ask me anything!")
			} else {
				// For simplicity, we'll just use a static API endpoint
				response, err := http.Get(openTriviaURL)
				if err != nil {
					log.Printf("Error fetching from API: %v", err)
					continue
				}
				defer response.Body.Close()

				var result map[string]interface{}
				err = json.NewDecoder(response.Body).Decode(&result)
				if err != nil {
					log.Printf("Error decoding API response: %v", err)
					continue
				}

				// Extract the question from the API response
				question := result["results"].([]interface{})[0].(map[string]interface{})["question"].(string)
				sendMessage(update.Message.Chat.ID, question)
			}
		}
	}
}

// getUpdates retrieves updates from Telegram using long polling
func getUpdates(offset int) ([]TelegramUpdate, error) {
	url := fmt.Sprintf(baseURL+getUpdatesURL, telegramBotToken, offset)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching updates: %w", err)
	}
	defer response.Body.Close()

	var responseData struct {
		Ok     bool             `json:"ok"`
		Result []TelegramUpdate `json:"result"`
	}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !responseData.Ok {
		return nil, fmt.Errorf("telegram API error")
	}

	return responseData.Result, nil
}

// sendMessage sends a message to a Telegram chat
func sendMessage(chatID int, text string) {
	message := TelegramResponse{
		ChatID: chatID,
		Text:   text,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return
	}

	// Send the message using Telegram Bot API
	_, err = http.Post(fmt.Sprintf(baseURL+"sendMessage", telegramBotToken),
		"application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}
