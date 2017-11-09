package botto

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type slackbot struct {
	id string

	client    *slack.Client
	rtm       *slack.RTM
	asyncEvts chan *interactiveResponse

	contexts map[string]*userContext
}

var log = logrus.WithField("app", "botto")

var bot *slackbot

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	if bot != nil {
		return errors.New("bot already connected")
	}

	bot = new(slackbot)
	bot.client = slack.New(token)
	bot.rtm = bot.client.NewRTM()
	bot.contexts = make(map[string]*userContext)
	bot.asyncEvts = make(chan *interactiveResponse)

	go bot.rtm.ManageConnection()

	for {
		select {
		case msg := <-bot.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				bot.id = ev.Info.User.ID

				log.Infof("%s is online @ %s", ev.Info.User.Name, ev.Info.Team.Name)

			case *slack.MessageEvent:
				bot.dispatchMsgEvent(ev)

			case *slack.RTMError:
				log.Errorf("RTM error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				return errors.New("invalid Slack token")
			}

		case ev := <-bot.asyncEvts:
			bot.dispatchAsync(ev)
		}
	}
}

func InteractiveEventHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(interactiveResponse)

	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bot.asyncEvts <- resp
}

func (b *slackbot) dispatchMsgEvent(ev *slack.MessageEvent) {
	// Only handle messages from other users
	if ev.User == b.id || (ev.Msg.Type != "message" ||
		(ev.Msg.SubType == "message_deleted" || ev.Msg.SubType == "bot_message")) {
		return
	}

	// Only handle direct messages
	if !strings.HasPrefix(ev.Msg.Channel, "D") {
		return
	}

	if _, ok := b.contexts[ev.User]; !ok {
		b.contexts[ev.User] = new(userContext)
		b.contexts[ev.User].init(&messageEvent{ev})
	}

	b.contexts[ev.User].dispatcher(b, &messageEvent{ev})
}

func (b *slackbot) dispatchAsync(resp *interactiveResponse) {
	ctx, ok := b.contexts[resp.User.ID]
	if !ok {
		log.Errorln("async event for non-existing user context")
		return
	}

	ctx.dispatcher(b, resp)
}
