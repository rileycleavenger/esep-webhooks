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

func FunctionHandler(ctx context.Context, event GitHubEvent) (string, error) {
    slackPayload := map[string]string{
        "text": fmt.Sprintf("Issue Created: %s", event.Issue.HTMLURL),
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