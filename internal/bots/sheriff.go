package bots

import (
	"time"

	"encoding/json"

	"fmt"

	"github.com/tbruyelle/hipchat-go/hipchat"

	"sort"

	"strings"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
)

type SheriffConfig struct {
	Time  string   `json:"time"`
	Topic string   `json:"topic"`
	Users []string `json:"users"`
}

type Sheriff struct {
	sheriffCfg     SheriffConfig
	wrangler       *core.Wrangler
	currentSheriff int
}

func (b *Sheriff) ticker() {
	ticker := time.NewTicker(1 * time.Second)
	fmtTime := fmt.Sprintf("%s:00", b.sheriffCfg.Time)
	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				//log15.Debug("msg", "now", now, "fmt", now.Format("15:04:05"), "cfg", fmtTime)
				if now.Format("15:04:05") == fmtTime {
					b.next()
				}
			}
		}
	}()
}

func (b *Sheriff) Name() string {
	return "Sheriff"
}

func (b *Sheriff) Description() string {
	return fmt.Sprintf("A bot that rotates users daily for the engineer on duty role a.k.a. sheriff. Current sheriff: %s Refresh time: %s", b.sheriffName(), b.sheriffCfg.Time)
}

func (b *Sheriff) Help() hipchat.NotificationRequest {
	help := `
	%s - %s
	<br>- /bot sheriff <strong>next</strong> - Switches to next sheriff
	<br>- /bot sheriff <strong>previous</strong> - Switches to previous sheriff
	<br>- /bot sheriff <strong>away</strong> $user - Marks the $user as away
	<br>- /bot sheriff <strong>back</strong> $user - Marks the $user as back
	`
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       fmt.Sprintf(help, b.Name(), b.Description()),
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b *Sheriff) extractAction(cmd string) string {
	w := strings.Fields(cmd)

	if len(w) >= 3 {
		// We can expect: /bot sheriff $action
		for _, v := range []string{"next", "previous", "away", "back"} {
			if w[2] == v {
				return v
			}
		}
	}

	return ""
}

func (b *Sheriff) HandleMessage(req *webhook.Request) hipchat.NotificationRequest {
	action := b.extractAction(req.Message())
	switch action {
	case "next":
		return b.next()
	case "previous":
		return b.previous()
	case "away":
		return b.status(true)
	case "back":
		return b.status(false)
	default:
		return b.unknown()
	}
}

func (b *Sheriff) HandleConfig(w *core.Wrangler, data json.RawMessage) error {
	var cfg SheriffConfig
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	// add config and boot
	b.sheriffCfg = cfg
	b.wrangler = w
	b.boot()

	return nil
}

func (b *Sheriff) boot() {
	// start ticker
	b.ticker()

	// sort sheriffs
	sort.Strings(b.sheriffCfg.Users)

	// pick first sheriff
	b.pickFirst()
}

func (b *Sheriff) next() hipchat.NotificationRequest {
	b.currentSheriff++
	if b.currentSheriff > len(b.sheriffCfg.Users) {
		// wrap around
		b.currentSheriff = 0
	}
	log15.Debug("switching to next sheriff", "sheriff", b.sheriffName())
	b.changeSheriff()
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       "We changed to the next sheriff: " + b.sheriffName(),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) previous() hipchat.NotificationRequest {
	b.currentSheriff--
	if b.currentSheriff < 0 {
		b.currentSheriff = len(b.sheriffCfg.Users) - 1
	}
	log15.Debug("switching to previous sheriff", "sheriff", b.sheriffName())
	b.changeSheriff()
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       "We switched back to the previous sheriff: " + b.sheriffName(),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) status(away bool) hipchat.NotificationRequest {
	status := "away"
	user := "you"
	log15.Debug("settings status of %s to %s", "user", user, "status", status)
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       "Marked $user as $status",
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) unknown() hipchat.NotificationRequest {
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorRed,
		Message:       "I'm sorry partner, don't know what you mean by that.. Run <strong>/bot sheriff --help</strong> to find out more.",
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b *Sheriff) changeSheriff() {
	if b.wrangler.Client != nil {
		go func() {
			topic := fmt.Sprintf(b.sheriffCfg.Topic, b.sheriffName())
			_, err := b.wrangler.Client.Room.SetTopic(b.wrangler.DefaultRoom, topic+" - "+time.Now().String())
			if err != nil {
				log15.Error("failed to set topic", "err", err)
			}
		}()
	}
}

func (b *Sheriff) pickFirst() {
	if len(b.sheriffCfg.Users) == 0 {
		// no sheriffs to pick from
		// TODO: This is hacky?
		b.currentSheriff = -1
		log15.Error("no sheriffs are configured")
		return
	}

	log15.Debug("picking first sheriff", "sheriff", b.sheriffName())
}

func (b *Sheriff) sheriffName() string {
	return b.sheriffCfg.Users[b.currentSheriff]
}
