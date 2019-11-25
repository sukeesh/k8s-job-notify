package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sukeesh/cron-k8s-watch/env"
)

type requestBody struct {
	Text string `json:"text"`
}

func SendSlackMessage(message string) error {
	slackBody, _ := json.Marshal(requestBody{Text: message})
	slackWebHookURL, err := env.GetSlackWebHookURL()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, slackWebHookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non ok response from Slack")
	}
	return nil
}
