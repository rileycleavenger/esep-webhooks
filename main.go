package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"io/ioutil"
)

type IssuePayload struct {
	Issue struct {
		HTMLURL string `json:"html_url"`
	} `json:"issue"`
}

func FunctionHandler(ctx context.Context, input []byte) (string, error) {
	var payload IssuePayload
	err := json.Unmarshal(input, &payload)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON input: %v", err)
	}

	slackPayload := map[string]string{
		"text": fmt.Sprintf("Issue Created: %s", payload.Issue.HTMLURL),
	}
	slackPayloadBytes, err := json.Marshal(slackPayload)
	if err != nil {
		return "", fmt.Errorf("error encoding Slack payload: %v", err)
	}

	slackURL := os.Getenv("SLACK_URL")
	if slackURL == "" {
		return "", fmt.Errorf("SLACK_URL environment variable is not set")
	}

	resp, err := http.Post(slackURL, "application/json", bytes.NewBuffer(slackPayloadBytes))
	if err != nil {
		return "", fmt.Errorf("error sending message to Slack: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from Slack: %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading Slack response body: %v", err)
	}

	return string(responseBody), nil
}

func main() {
	input := []byte(`{"issue": {"html_url": "https://example.com/issue"}}`)
	result, err := FunctionHandler(context.Background(), input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Response from Slack:", result)
}
