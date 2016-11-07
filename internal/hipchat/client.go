package hipchat

import (
	"bytes"
	"net/http"
	"strconv"

	"encoding/json"
	"fmt"
)

type Client struct {
	Url         string `json:"url"`
	Auth        string `json:"auth"`
	DefaultRoom int    `json:"roomId"`
}

func NewClient(url, auth string, room int) *Client {
	return &Client{
		Url:         url,
		Auth:        auth,
		DefaultRoom: room,
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
	url := c.Url + "/room/" + strconv.Itoa(c.DefaultRoom) + "/notification?auth_token=" + c.Auth

	nb, _ := json.Marshal(notification)

	return c.send("POST", url, nb)
}

type topicResponse struct {
	string `json:"topic"`
}

func (c *Client) ChangeRoomTopic(topic string) error {
	url := fmt.Sprintf("%s/room/%d/topic?auth_token=%s", c.Url, c.DefaultRoom, c.Auth)

	tp := topicResponse{topic}
	body, _ := json.Marshal(tp)
	return c.send("PUT", url, []byte(body))
}
