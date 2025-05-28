package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	geminiAPIURLBase   = "https://generativelanguage.googleapis.com/v1beta/models/"
	geminiModelDefault = "gemini-2.0-flash-thinking-exp-01-21"
	geminiAction       = ":generateContent?key="
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ContentResponse `json:"content"`
}

type ContentResponse struct {
	Parts []PartResponse `json:"parts"`
}

type PartResponse struct {
	Text string `json:"text"`
}

func Summarize(userInput string, systemPrompt string, apiKey string) (string, error) {
	url := fmt.Sprintf("%s%s%s%s", geminiAPIURLBase, geminiModelDefault, geminiAction, apiKey)

	fullPrompt := userInput
	if systemPrompt != "" {
		fullPrompt = fmt.Sprintf("%s\n\n%s", systemPrompt, userInput)
	}

	requestBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: fullPrompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error serializing request body to JSON: %w", err) 
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", fmt.Errorf("error deserializing response JSON: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates found in response")
	}
	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts found in candidate's content")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
