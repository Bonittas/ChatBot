package openapi

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "log"
)

// Client represents the API client to interact with the OpenAPI service.
type Client struct {
    apiUrl string
}

// NewClient creates a new instance of Client with the provided API URL.
func NewClient(apiUrl string) *Client {
    return &Client{apiUrl: apiUrl}
}

// GetAnswer makes a request to the OpenAPI service to get an answer to the given question.
func (c *Client) GetAnswer(question string) (string, error) {
    apiUrl, err := url.Parse(c.apiUrl)
    if err != nil {
        return "", err
    }

    apiUrl.Path += "/answer"
    query := apiUrl.Query()
    query.Set("question", question)
    apiUrl.RawQuery = query.Encode()

    log.Printf("Requesting URL: %s", apiUrl.String())
    resp, err := http.Get(apiUrl.String())
    if err != nil {
        log.Printf("HTTP request failed: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    log.Printf("Received response with status code: %d", resp.StatusCode)
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var result AnswerResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        log.Printf("Failed to decode response: %v", err)
        return "", err
    }

    log.Printf("Received answer: %s", result.Answer)
    return result.Answer, nil
}
