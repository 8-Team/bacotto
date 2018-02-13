package slackbot

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type Bot struct {
	id string

	client *slack.Client
	rtm    *slack.RTM

	ctxInit  ContextFn
	contexts sync.Map
	commands sync.Map

	logger *logrus.Logger
}

type command struct {
	name        string
	description string
	callback    ContextFn
}

// New creates a new Slack bot using the provided API token.
func New(token string) *Bot {
	bot := new(Bot)

	bot.client = slack.New(token)
	bot.rtm = bot.client.NewRTM()

	bot.logger = logrus.New()
	bot.logger.SetLevel(logrus.DebugLevel)

	return bot
}

// Start starts receiving events and triggers the appropriate registered context callbacks.
func (bot *Bot) Start() error {
	go bot.rtm.ManageConnection()

	for {
		msg := <-bot.rtm.IncomingEvents

		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			bot.logger.Infof("%s is online @ %s", ev.Info.User.Name, ev.Info.Team.Name)
			bot.id = ev.Info.User.ID

		case *slack.MessageEvent:
			bot.logger.Debugln("Message event received")
			bot.dispatchMsgEvent(ev)

		case *slack.RTMError:
			bot.logger.Errorf("RTM error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			return errors.New("invalid Slack token")
		}
	}
}

// InteractiveEventHandler listens to incoming interactive (asynchronous) messages and
// dispatches them to the correct context.
func (bot *Bot) InteractiveEventHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(slack.AttachmentActionCallback)
	payload := r.FormValue("payload")

	if err := json.Unmarshal([]byte(payload), resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bot.logger.Errorln(err)
		return
	}

	// Receiving an async message from a non-existing user is a no-no.
	ctx, ok := bot.contexts.Load(resp.User.ID)
	if !ok {
		bot.logger.Errorln("Asynchronous event for non-existing user, skipping")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.(*Context).dispatchAsync(resp)
	w.WriteHeader(http.StatusOK)
}

// dispatchMsgEvent handles incoming "synchronous" messages from the bot RTM API
func (bot *Bot) dispatchMsgEvent(ev *slack.MessageEvent) {
	// Only handle messages from other users
	if ev.User == "" || ev.User == bot.id || (ev.Msg.Type != "message" ||
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
		bot.logger.Debugf("Missing user context for %s, creating one", ev.User)
		go ctx.(*Context).start(bot.ctxInit)
	}

	ctx.(*Context).dispatchSync(ev)
}

// OnContextCreation will set a context creation function invoked as soon as a new context
// is scheduled for creation. If an error is returned, the context creation is aborted and
// the error message is shown to the user.
func (bot *Bot) OnContextCreation(initFn ContextFn) {
	bot.ctxInit = initFn
}

// Register associates a callback to a command. The callback is invoked whenever a message
// is received by the context matching the text.
func (bot *Bot) Register(text, description string, f ContextFn) {
	bot.commands.Store(text, command{
		name:        text,
		description: description,
		callback:    f,
	})
}
