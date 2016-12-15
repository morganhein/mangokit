package discord

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
	"github.com/morganhein/mangokit/events"
)

// event handlers
func onConnect(s *discordgo.Session, _ *discordgo.Connect) {
	g, _ := s.Gateway()
	event := plugins.Event{
		Time: time.Now(),
		Context: &guild{s: s},
		Type: events.CONNECTED,
		Raw: "Connected to " + g,
	}
	toApp <- event
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == server.me.ID {
		return
	}
	log.Debug("Received new message from: " + m.Author.Username + ":" + m.Author.ID)

	cid, err := s.Channel(m.ChannelID)

	if err != nil {
		//todo: maybe some better checking here?
		cid, _ = s.Channel(m.Message.ChannelID)
	}

	if (cid == nil) {
		log.Warning("Unable to determine channel ID for incoming message.")
	}

	channel := &channel{
		s: s,
		c: cid,
	}

	event := plugins.Event{
		Time: time.Now(),
		Context: channel,
		Type: events.MESSAGE,
		Raw: m.Message.Content,
		Who: &plugins.Who{
			Name: m.Author.Username,
			Id: m.Author.ID,
		},
	}
	toApp <- event
}
