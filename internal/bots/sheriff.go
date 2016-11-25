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
		return w[2]
	}

	return ""
}

func (b *Sheriff) HandleMessage(req *webhook.Request) (r hipchat.NotificationRequest, useNotification bool) {
	useNotification = false
	action := b.extractAction(req.Message())
	switch action {
	case "next":
		return b.next(), true
	case "previous":
		return b.previous(), true
	case "away":
		return b.status(req, true), true
	case "back":
		return b.status(req, false), true
	case "list":
		return b.list(), true
	case "test":
		go func() {
			n := hipchat.NotificationRequest{
				Color:         hipchat.ColorPurple,
				Message:       "We are sending a test notification! Thanks, " + req.Username(),
				Notify:        false,
				MessageFormat: "text",
			}
			b.wrangler.SendNotification(b, &n)
		}()
		return
	default:
		return b.unknown(), true
	}
}

func (b *Sheriff) HandleConfig(w *core.Wrangler, data json.RawMessage) error {
	var cfg SheriffConfig
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	// do we have sheriffs?
	if len(cfg.Users) == 0 {
		return fmt.Errorf("you have not supplied sheriffs in the 'users' field")
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
}

func (b *Sheriff) next() hipchat.NotificationRequest {
	b.currentSheriff++
	if b.currentSheriff > len(b.sheriffCfg.Users) {
		// wrap around
		b.currentSheriff = 0
	}
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
	b.changeSheriff()
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       "We switched back to the previous sheriff: " + b.sheriffName(),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) status(req *webhook.Request, away bool) hipchat.NotificationRequest {
	user := req.Username()
	status := ""
	msg := ""
	if away {
		status = "away"
		msg = "Bye! :-*"
	} else {
		status = "back"
		msg = "Welcome back! :-D"
	}
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       fmt.Sprintf("Marking user %s as %s. %s", user, status, msg),
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

func (b *Sheriff) list() hipchat.NotificationRequest {
	sheriffs := strings.Join(b.sheriffCfg.Users, "<br>")

	return hipchat.NotificationRequest{
		Color:         hipchat.ColorGray,
		Message:       fmt.Sprintf("These are the sheriffs of the town:<br>%s", sheriffs),
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b *Sheriff) changeSheriff() {
	if b.wrangler.Client != nil {
		go func() {
			topic := fmt.Sprintf(b.sheriffCfg.Topic, b.sheriffName())
			_, err := b.wrangler.Client.Room.SetTopic(b.wrangler.DefaultRoom, topic)
			if err != nil {
				log15.Error("failed to set topic", "err", err)
			}
		}()
	}
}

func (b *Sheriff) sheriffName() string {
	return b.sheriffCfg.Users[b.currentSheriff]
}
