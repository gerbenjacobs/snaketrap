package bots

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

type EightBall struct {
	wrangler *core.Wrangler
}

func (b *EightBall) Name() string {
	return "8ball"
}

func (b *EightBall) Description() string {
	return "C"
}

func (b *EightBall) Help() core.Reply {
	help := `
	%s - %s
	<br>- /bot 8ball $question - The 8-ball will answer your question!
	`
	return core.NewReply(hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       fmt.Sprintf(help, b.Name(), b.Description()),
		Notify:        false,
		MessageFormat: "html",
	})
}

func (b *EightBall) HandleMessage(req *webhook.Request) core.Reply {
	question := req.Message()

	go b.reply()

	return core.NewReply(hipchat.NotificationRequest{
		Color:         hipchat.ColorGray,
		Message:       fmt.Sprintf("You asked: %s - Let me think about that..", question),
		Notify:        false,
		MessageFormat: "text",
	})
}

func (b *EightBall) HandleConfig(w *core.Wrangler, data json.RawMessage) error {
	b.wrangler = w
	return nil
}

func (b *EightBall) CurrentState() []byte {
	return nil
}

func (b *EightBall) reply() {
	time.Sleep(2 * time.Second)

	rand.Seed(time.Now().Unix())
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}

	n := hipchat.NotificationRequest{
		Color:         hipchat.ColorPurple,
		Message:       fmt.Sprintf(".. my answer is: %s", answers[rand.Intn(len(answers))]),
		Notify:        false,
		MessageFormat: "text",
	}

	log15.Debug("replying to previous 8ball request", "n", n)
	b.wrangler.SendNotification(b, &n)
}
