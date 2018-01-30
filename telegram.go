// Package telegram provides a very basic and simple mean to send messages
// to a specific target (user/channel/group) via a chatbot.
package telegram_msg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const sendMessageURLTpl = "https://api.telegram.org/bot%s/sendMessage"

// Telegram has all the required dynamic information.
type Telegram struct {
	client   *http.Client
	botToken string
	targetID int
}

type telegramMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	Silent bool   `json:"disable_notification"`
}

// Create initiates the Telegram Bot messenger using all the vital information provided.
func (s *Telegram) Create(botToken string, targetID int) {
	s.botToken = botToken
	s.targetID = targetID
	s.client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

// Send uses the Stringer interface to compose the messages sent over to Slack.
func (s *Telegram) Send(text string, params ...interface{}) {
	send(s, fmt.Sprintf(text, params...), false)
}

func (s *Telegram) SendSilent(text string, params ...interface{}) {
	send(s, fmt.Sprintf(text, params...), true)
}

func send(telegram *Telegram, text string, silent bool) {
	t := &telegramMessage{
		ChatId: telegram.targetID,
		Text:   text,
		Silent: silent}
	payload, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("Failed to create json payload for Telegram Bot: %s\n",
			err.Error())
	}

	p := strings.NewReader(string(payload))
	resp, err := telegram.client.Post(
		fmt.Sprintf(sendMessageURLTpl, telegram.botToken),
		"application/json", p)
	if err != nil {
		log.Fatalf("Failed to pass text to Telegram Bot: %s\n", err.Error())
	}
	defer resp.Body.Close()
}
