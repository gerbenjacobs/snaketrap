package bots

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

const (
	EmojiUnauthorized = "(foreveralone)"
	EmojiFailed       = "(failed)"
	EmojiSuccess      = "(successful)"

	CmdNext     = "next"
	CmdPrevious = "previous"
	CmdSet      = "set"
	CmdAway     = "away"
	CmdBack     = "back"
	CmdList     = "list"
	CmdAnnounce = "announce"
)

var AuthorizedCmds = []string{CmdNext, CmdPrevious, CmdSet, CmdAway, CmdBack}

type SheriffConfig struct {
	Days         []int        `json:"days"`
	Time         string       `json:"time"`
	Announce     bool         `json:"announce"`
	AnnounceTime string       `json:"announce_time"`
	Topic        string       `json:"topic"`
	Users        []string     `json:"users"`
	SheriffUsers SheriffUsers `json:"-"`
}

type Sheriff struct {
	config   SheriffConfig
	wrangler *core.Wrangler
	current  int
}

type SheriffUsers []SheriffUser

type SheriffUser struct {
	Name string `json:"name"`
	Away bool   `json:"away"`
}

type SheriffState struct {
	CurrentSheriff int          `json:"current_sheriff"`
	SheriffUsers   SheriffUsers `json:"sheriff_users"`
}

func (b *Sheriff) ticker() {
	ticker := time.NewTicker(1 * time.Second)
	rotateTime := fmt.Sprintf("%s:00", b.config.Time)
	announceTime := fmt.Sprintf("%s:00", b.config.AnnounceTime)
	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				if b.isActiveDay(now) && now.Format("15:04:05") == rotateTime {
					n := b.rotate(true)
					b.wrangler.SendNotification(b, &n)
					log15.Debug("switching to new sheriff", "sherrif", b.sheriffName(), "time", now)
				}

				if b.config.Announce && now.Format("15:04:05") == announceTime {
					go b.sendAnnouncement()
				}
			}
		}
	}()
}

func (b *Sheriff) Name() string {
	return "Sheriff"
}

func (b *Sheriff) Description() string {
	return fmt.Sprintf("A bot that rotates users daily for the engineer on duty role a.k.a. sheriff. Current sheriff: %s Refresh time: %s", b.sheriffName(), b.config.Time)
}

func (b *Sheriff) Help() hipchat.NotificationRequest {
	help := fmt.Sprintf(`
	<strong>%s</strong> - %s
	<br>- /bot sheriff <strong>%s</strong> - * Switches to next sheriff
	<br>- /bot sheriff <strong>%s</strong> - * Switches to previous sheriff
	<br>- /bot sheriff <strong>%s</strong> $user - * Switches sheriff duties to $user
	<br>- /bot sheriff <strong>%s</strong> $user - * Marks the $user as away
	<br>- /bot sheriff <strong>%s</strong> $user - * Marks the $user as back
	<br>- /bot sheriff <strong>%s</strong> - Lists the sheriffs and their availability
	<br>- /bot sheriff <strong>%s</strong> - Announces the sheriff for the next day
	<br><br><em>* admin only</em>
	`,
		b.Name(), b.Description(),
		CmdNext, CmdPrevious, CmdSet, CmdAway, CmdBack, CmdList, CmdAnnounce,
	)

	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       help,
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

func (b *Sheriff) HandleMessage(req *webhook.Request) (hipchat.NotificationRequest, bool) {
	action := b.extractAction(req.Message())
	requester := req.Username()

	// enforce authorization
	for _, a := range AuthorizedCmds {
		if a == action && !b.isAuthorized(requester) {
			return b.unauthorized(), true
		}
	}

	switch action {
	case CmdNext:
		return b.rotate(true), true
	case CmdPrevious:
		return b.rotate(false), true
	case CmdSet:
		return hipchat.NotificationRequest{
			Color:   hipchat.ColorYellow,
			Message: "This method is not implemented yet..",
		}, true
	case CmdAway:
		return b.status(req, true), true
	case CmdBack:
		return b.status(req, false), true
	case CmdList:
		return b.list(), true
	case CmdAnnounce:
		go b.sendAnnouncement()
		return hipchat.NotificationRequest{}, false
	case "test":
		go func() {
			n := hipchat.NotificationRequest{
				Color:         hipchat.ColorPurple,
				Message:       "We are sending a test notification! Thanks, " + requester,
				Notify:        false,
				MessageFormat: "text",
			}
			b.wrangler.SendNotification(b, &n)
		}()
		return hipchat.NotificationRequest{}, false
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
		return errors.New("you have not supplied sheriffs in the 'users' field")
	}

	// add config and boot
	b.config = cfg
	b.wrangler = w
	b.boot()

	// upsert config in state file
	err = b.initializeState()
	if err != nil {
		log15.Error("failed to initialize config in state file", "err", err)
	}

	return nil
}

func (b *Sheriff) CurrentState() []byte {
	// create state
	s := SheriffState{
		CurrentSheriff: b.current,
		SheriffUsers:   b.config.SheriffUsers,
	}

	// marshal to json
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log15.Error("failed to marshal config", "err", err)
		return nil
	}

	return data
}

func (b *Sheriff) boot() {
	// start ticker
	b.ticker()

	// create SheriffUsers
	for _, u := range b.config.Users {
		b.config.SheriffUsers = append(b.config.SheriffUsers, SheriffUser{
			Name: u,
			Away: false,
		})
	}

	// sort sheriffs
	sort.Sort(b.config.SheriffUsers)
}

func (b *Sheriff) nextSheriff(next bool, current int) int {
	if next {
		current++
		if current >= len(b.config.SheriffUsers) {
			// wrap around
			current = 0
		}
	} else {
		current--
		if current < 0 {
			current = len(b.config.SheriffUsers) - 1
		}
	}

	if b.config.SheriffUsers[current].Away {
		// sheriff is unavailable, rotate again
		log15.Debug("skipping unavailable sheriff", "sherrif", b.sheriffNameById(current))
		return b.nextSheriff(next, current)
	}

	return current
}

func (b *Sheriff) rotate(next bool) hipchat.NotificationRequest {
	b.current = b.nextSheriff(next, b.current)
	go b.changeSheriff()
	go b.setState()

	var msg string
	if next {
		msg = fmt.Sprintf("We changed to the next sheriff: %s", b.sheriffName())
	} else {
		msg = fmt.Sprintf("We switched back to the previous sheriff: %s", b.sheriffName())
	}

	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       msg,
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) status(req *webhook.Request, away bool) hipchat.NotificationRequest {
	user := req.GetWord(3)
	if user == "" {
		user = req.Username()
	}

	// store status
	found := false
	for i, u := range b.config.SheriffUsers {
		if u.Name == user {
			found = true
			b.config.SheriffUsers[i].Away = away
		}
	}

	if !found {
		return hipchat.NotificationRequest{
			Color:         hipchat.ColorRed,
			Message:       fmt.Sprintf("We can't seem to find the sheriff named %s", user),
			Notify:        false,
			MessageFormat: "text",
		}
	}

	b.setState()

	// send reply
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
		Message:       fmt.Sprintf("Marking user %s as %s. %s", req.Fullname(), status, msg),
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

func (b *Sheriff) unauthorized() hipchat.NotificationRequest {
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorRed,
		Message:       "I'm sorry partner, but you don't seem to be holding a shiny gold sheriff's badge! " + EmojiUnauthorized,
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b *Sheriff) sendAnnouncement() {
	if b.isActiveDayTomorrow(time.Now()) {
		rotatedSheriff := b.nextSheriff(true, b.current)
		s := b.sheriffNameById(rotatedSheriff)
		n := hipchat.NotificationRequest{
			Color:   hipchat.ColorPurple,
			Message: "Tomorrow's sheriff will be: " + s,
		}
		b.wrangler.SendNotification(b, &n)
		log15.Debug("announcing sheriff for the next day", "sheriff", s, "time", time.Now())
	} else {
		n := hipchat.NotificationRequest{
			Color:   hipchat.ColorPurple,
			Message: "There is no sheriff's duty tomorrow..",
		}
		b.wrangler.SendNotification(b, &n)
		log15.Debug("no sheriff duty tomorrow..")
	}
}

func (b *Sheriff) list() hipchat.NotificationRequest {
	n := []string{}
	for _, u := range b.config.SheriffUsers {
		n = append(n, u.sheriffStatus()+" "+u.Name)
	}
	sheriffs := strings.Join(n, "\n")

	return hipchat.NotificationRequest{
		Color:         hipchat.ColorGray,
		Message:       fmt.Sprintf("These are the sheriffs of the town:\n%s", sheriffs),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (slice SheriffUsers) Len() int {
	return len(slice)
}

func (slice SheriffUsers) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice SheriffUsers) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (b *Sheriff) changeSheriff() {
	topic := fmt.Sprintf(b.config.Topic, b.sheriffName())
	b.wrangler.SetTopic(b, topic)
}

func (b *Sheriff) sheriffName() string {
	return b.config.SheriffUsers[b.current].Name
}

func (b *Sheriff) sheriffNameById(id int) string {
	return b.config.SheriffUsers[id].Name
}

func (u SheriffUser) sheriffStatus() string {
	if u.Away {
		return EmojiFailed
	}
	return EmojiSuccess
}

func (b *Sheriff) isActiveDay(now time.Time) bool {
	for _, d := range b.config.Days {
		if int(now.Weekday()) == d {
			return true
		}
	}
	return false
}

func (b *Sheriff) isActiveDayTomorrow(now time.Time) bool {
	return b.isActiveDay(now.AddDate(0, 0, 1))
}

func (b *Sheriff) setState() {
	go func() {
		if err := b.wrangler.SetState(b); err != nil {
			log15.Error("failed to set state", "err", err)
		}
	}()
}

func (b *Sheriff) initializeState() error {
	// retrieve
	d, err := b.wrangler.GetState(b)
	if err != nil {
		return err
	}

	// unmarshal into state
	var s SheriffState
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}

	// apply state
	b.current = s.CurrentSheriff
	b.mergeSheriffUsers(s.SheriffUsers)

	return nil
}

func (b *Sheriff) mergeSheriffUsers(state SheriffUsers) {
	status := map[string]bool{}
	for _, stateSheriff := range state {
		status[stateSheriff.Name] = stateSheriff.Away
	}

	for i, configSheriff := range b.config.SheriffUsers {
		if availability, ok := status[configSheriff.Name]; ok {
			b.config.SheriffUsers[i].Away = availability
		}
	}
}

func (b *Sheriff) isAuthorized(sheriff string) bool {
	for _, s := range b.config.SheriffUsers {
		if s.Name == sheriff {
			return true
		}
	}

	return false
}
