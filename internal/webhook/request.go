package webhook

import "strings"

type Request struct {
	Event     string `json:"event"`
	Item      Item   `json:"item"`
	WebhookID int    `json:"webhook_id"`
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
	a := strings.Fields(r.Item.Message.Message)
	if len(a) < 2 {
		return ""
	}
	return strings.Join(a[2:], " ")
}

func (r Request) Username() string {
	return r.Item.Message.From.MentionName
}

func (r Request) Fullname() string {
	return r.Item.Message.From.Name
}

func (r Request) GetToken(n int) string {
	w := strings.Fields(r.Item.Message.Message)
	if n >= len(w) {
		return ""
	}
	return w[n]
}

func (r Request) Bot() string {
	return r.GetToken(1)
}
