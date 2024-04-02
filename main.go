package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "github.com/aws/aws-lambda-go/lambda"
)

type GitHubEvent struct {
    Issue struct {
        HTMLURL string `json:"html_url"`
    } `json:"issue"`
}

type Event struct {
    Body string `json:"body"`
}

func FunctionHandler(ctx context.Context, input json.RawMessage) (string, error) {
    var e Event
    if err := json.Unmarshal(input, &e); err != nil {
        return "", fmt.Errorf("error decoding input: %v", err)
    }

    var gitHubEvent GitHubEvent
    if err := json.Unmarshal([]byte(e.Body), &gitHubEvent); err != nil {
        return "", fmt.Errorf("error decoding GitHub event: %v", err)
    }

    slackPayload := map[string]string{
        "text": fmt.Sprintf("Issue Created: %s", gitHubEvent.Issue.HTMLURL),
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
    lambda.Start(FunctionHandler)
}