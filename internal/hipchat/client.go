package hipchat

import (
	"bytes"
	"net/http"
	"strconv"

	"encoding/json"
	"fmt"
)

type Client struct {
	url          string
	auth         string
	default_room int
}

func NewClient(url, auth string, room int) *Client {
	return &Client{
		url:          url,
		auth:         auth,
		default_room: room,
	}
}

func (c *Client) send(method, url string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendNotification(notification *Response) error {
	url := c.url + "/room/" + strconv.Itoa(c.default_room) + "/notification?auth_token=" + c.auth

	nb, _ := json.Marshal(notification)

	return c.send("POST", url, nb)
}

type topicResponse struct {
	string `json:"topic"`
}

func (c *Client) ChangeRoomTopic(topic string) error {
	url := fmt.Sprintf("%s/room/%d/topic?auth_token=%s", c.url, c.default_room, c.auth)

	tp := topicResponse{topic}
	body, _ := json.Marshal(tp)
	return c.send("PUT", url, []byte(body))
}
