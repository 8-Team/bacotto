package botto

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/8-team/bacotto/conf"
	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type slackbot struct {
	id string

	client    *slack.Client
	rtm       *slack.RTM
	asyncEvts chan *interactiveResponse
	callbacks map[string]eventCallback

	contexts map[string]*userContext
}

type eventCallback func(*slackbot, *interactiveResponse)

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
	bot.contexts = make(map[string]*userContext)
	bot.asyncEvts = make(chan *interactiveResponse)
	bot.callbacks = make(map[string]eventCallback)

	go bot.rtm.ManageConnection()

	for {
		select {
		case msg := <-bot.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				bot.id = ev.Info.User.ID

				log.Infof("%s is online @ %s", ev.Info.User.Name, ev.Info.Team.Name)

			case *slack.MessageEvent:
				log.Debugln("Message event received")
				bot.dispatchMsgEvent(ev)

			case *slack.RTMError:
				log.Errorf("RTM error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				return errors.New("invalid Slack token")
			}

		case ev := <-bot.asyncEvts:
			log.Debugln("Async event received")
			bot.dispatchAsync(ev)
		}
	}
}

func InteractiveEventHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(interactiveResponse)
	payload := r.FormValue("payload")

	if err := json.Unmarshal([]byte(payload), resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	bot.asyncEvts <- resp

	w.WriteHeader(http.StatusOK)
}

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

	if _, ok := b.contexts[ev.User]; !ok {
		log.Debugf("Missing user context for %s, creating one", ev.User)
		b.contexts[ev.User] = new(userContext)
		b.contexts[ev.User].init(&messageEvent{ev})
	}

	b.contexts[ev.User].dispatcher(b, &messageEvent{ev})
}

func (b *slackbot) dispatchAsync(resp *interactiveResponse) {
	eventCallback, ok := b.callbacks[resp.CallbackID]
	if !ok {
		log.Errorln("invalid callback for async response")
		return
	}

	if eventCallback != nil {
		eventCallback(b, resp)
		b.removeCallback(resp.CallbackID)
	}
}

func (b *slackbot) registerCallback(id string, cb eventCallback) {
	b.callbacks[id] = cb
}

func (b *slackbot) removeCallback(id string) {
	delete(b.callbacks, id)
}
