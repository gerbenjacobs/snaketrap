package webhook

import "strings"

type Request struct {
	Event     string `json:"event"`
	Item      Item   `json:"item"`
	WebhookId int    `json:"webhook_id"`
}

type Item struct {
	Message Message `json:"message"`
	Room    Room    `json:"room"`
}

type Message struct {
	ID       string   `json:"id"`
	Date     string   `json:"date"`
	From     From     `json:"from"`
	Message  string   `json:"message"`
	Type     string   `json:"type"`
	Mentions []string `json:"mentions"`
}

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type From struct {
	ID          int    `json:"id"`
	MentionName string `json:"mention_name"`
	Name        string `json:"name"`
}

func (r Request) Message() string {
	return r.Item.Message.Message
}

func (r Request) Username() string {
	return r.Item.Message.From.MentionName
}

func (r Request) GetWord(n int) string {
	w := strings.Fields(r.Message())
	if n >= len(w) {
		return ""
	}
	return w[n]
}

func (r Request) Bot() string {
	return r.GetWord(1)
}
