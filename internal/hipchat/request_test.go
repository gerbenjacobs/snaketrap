package hipchat

import (
	"encoding/json"
	"testing"
)

var jsData = `
{
	"event": "room_message",
	"item": {
		"message": {
			"date": "2016-11-04T12:38:24.363690+00:00",
			"from": {
				"id": 1073020,
				"links": {
					"self": "https://api.hipchat.com/v2/user/1073020"
				},
				"mention_name": "gerben",
				"name": "Gerben Jacobs",
				"version": "00000000"
			},
			"id": "b4fb28bd-617b-4661-a1cd-ecefb83680f0",
			"mentions": [],
			"message": "/bot my message here",
			"type": "message"
		},
		"room": {
			"id": 3277014,
			"is_archived": false,
			"links": {
				"participants": "https://api.hipchat.com/v2/room/3277014/participant",
				"self": "https://api.hipchat.com/v2/room/3277014",
				"webhooks": "https://api.hipchat.com/v2/room/3277014/webhook"
			},
			"name": "Snaketrap",
			"privacy": "public",
			"version": "IV7L53Z0"
		}
	},
	"oauth_client_id": "5d21ad94-9b26-41ba-9a14-539d640b5394",
	"webhook_id": 10519583
}`

func TestMarshalling(t *testing.T) {
	req := new(Request)
	err := json.Unmarshal([]byte(jsData), req)

	if err != nil {
		t.Error("Failed to unmarshal hipchat.Request", "err", err)
	}
}
