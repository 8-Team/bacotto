package botto

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/8-team/bacotto/conf"
	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type slackbot struct {
	id string

	client   *slack.Client
	rtm      *slack.RTM
	contexts sync.Map
}

var log *logrus.Entry
var bot *slackbot

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	if bot != nil {
		return errors.New("bot already connected")
	}

	if conf.DebugLogLevel() {
		logrus.SetLevel(logrus.DebugLevel)
	}
	log = logrus.WithField("app", "botto")

	bot = new(slackbot)
	bot.client = slack.New(token)
	bot.rtm = bot.client.NewRTM()

	go bot.rtm.ManageConnection()

	for {
		msg := <-bot.rtm.IncomingEvents

		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Infof("%s is online @ %s", ev.Info.User.Name, ev.Info.Team.Name)
			bot.id = ev.Info.User.ID

		case *slack.MessageEvent:
			log.Debugln("Message event received")
			bot.dispatchMsgEvent(ev)

		case *slack.RTMError:
			log.Errorf("RTM error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			return errors.New("invalid Slack token")
		}
	}
}

// InteractiveEventHandler listens to incoming interactive (asynchronous) messages and
// dispatches them to the correct context.
func InteractiveEventHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(slack.AttachmentActionCallback)
	payload := r.FormValue("payload")

	if err := json.Unmarshal([]byte(payload), resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}

	// Receiving an async message from a non-existing user is a no-no.
	ctx, ok := bot.contexts.Load(resp.User.ID)
	if !ok {
		log.Errorln("Asynchronous event for non-existing user, skipping")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.(*context).dispatchAsync(resp)
	w.WriteHeader(http.StatusOK)
}

// dispatchMsgEvent handles incoming "synchronous" messages from the bot RTM API
func (b *slackbot) dispatchMsgEvent(ev *slack.MessageEvent) {
	// Only handle messages from other users
	if ev.User == "" || ev.User == b.id || (ev.Msg.Type != "message" ||
		(ev.Msg.SubType == "message_deleted" || ev.Msg.SubType == "bot_message")) {
		return
	}

	// Only handle direct messages
	if !strings.HasPrefix(ev.Msg.Channel, "D") {
		return
	}

	// A user context is created and started if this is the first message received from that user.
	ctx, loaded := bot.contexts.LoadOrStore(ev.User, newContext(bot, ev.User, ev.Channel))
	if !loaded {
		log.Debugf("Missing user context for %s, creating one", ev.User)
		go ctx.(*context).start()
	}

	ctx.(*context).dispatchSync(ev)
}
