package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/plugins"
)

// event handlers
func (d *discord) onConnect(s *discordgo.Session, _ *discordgo.Connect) {
	g, _ := s.Gateway()
	event := plugins.Event{
		Time:    time.Now(),
		Context: &guild{s: s},
		Type:    events.CONNECTED,
		Raw:     "Connected to " + g,
	}
	disc.con.FromPlugin <- event
}

func (d *discord) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == disc.me.ID {
		return
	}
	log.Debug("Received new message from: " + m.Author.Username + ":" + m.Author.ID)

	cid, err := s.Channel(m.ChannelID)

	if err != nil {
		//todo: maybe some better checking here?
		cid, _ = s.Channel(m.Message.ChannelID)
	}

	if cid == nil {
		log.Warning("Unable to determine channel ID for incoming message.")
	}

	channel := &channel{
		s: s,
		c: cid,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.MESSAGE,
		Raw:     m.Message.Content,
		Who: &plugins.Who{
			Name: m.Author.Username,
			Id:   m.Author.ID,
		},
	}
	disc.con.FromPlugin <- event
}

func (d *discord) onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == disc.me.ID {
		return
	}

	log.Debug("Received updated message from: " + m.Author.Username + ":" + m.Author.ID)

	cid, err := s.Channel(m.ChannelID)

	if err != nil {
		//todo: maybe some better checking here?
		cid, _ = s.Channel(m.Message.ChannelID)
	}

	if cid == nil {
		log.Warning("Unable to determine channel ID for updated message.")
	}

	channel := &channel{
		s: s,
		c: cid,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.UPDATEMESSAGE,
		Raw:     m.Message.Content,
		Who: &plugins.Who{
			Name: m.Author.Username,
			Id:   m.Author.ID,
		},
	}
	disc.con.FromPlugin <- event
}

func (d *discord) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == disc.me.ID {
		return
	}

	log.Debug("Deleted message from: " + m.Author.Username + ":" + m.Author.ID)

	ch, err := s.Channel(m.ChannelID)

	if err != nil {
		//todo: maybe some better checking here?
		ch, _ = s.Channel(m.Message.ChannelID)
	}

	if ch == nil {
		log.Warning("Unable to determine channel ID for deleted message.")
	}

	channel := &channel{
		s: s,
		c: ch,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.DESTROYMESSAGE,
		Who: &plugins.Who{
			Name: m.Author.Username,
			Id:   m.Author.ID,
		},
	}
	disc.con.FromPlugin <- event
}

func (d *discord) onChannelJoin(s *discordgo.Session, c *discordgo.ChannelCreate) {
	channel := &channel{
		s: s,
		c: c.Channel,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.CREATEDCHANNEL,
	}

	disc.con.FromPlugin <- event
}

func (d *discord) onGuildMemberAdd(s *discordgo.Session, m *discordgo.Member) {
	// Ignore all events created by the bot itself
	if m.User.ID == d.me.ID {
		return
	}

	// get the welcome channel for that guild
	ch, err := s.Channel(m.GuildID)
	if err != nil {
		return
	}

	if ch == nil {
		log.Warning("Unable to determine default channel ID. This event cannot be handled.")
	}

	channel := &channel{
		s: s,
		c: ch,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.JOINEDSERVER,
		Who: &plugins.Who{
			Name: m.User.Username,
			Id:   m.User.ID,
		},
	}
	disc.con.FromPlugin <- event
}

func (d *discord) onGuildMemberExit(s *discordgo.Session, m *discordgo.Member) {
	// Ignore all events created by the bot itself
	if m.User.ID == d.me.ID {
		return
	}

	// get the welcome channel for that guild
	ch, err := s.Channel(m.GuildID)
	if err != nil {
		return
	}

	if ch == nil {
		log.Warning("Unable to determine default channel ID. This event cannot be handled.")
	}

	channel := &channel{
		s: s,
		c: ch,
	}

	event := plugins.Event{
		Time:    time.Now(),
		Context: channel,
		Type:    events.EXITEDSERVER,
		Who: &plugins.Who{
			Name: m.User.Username,
			Id:   m.User.ID,
		},
	}
	disc.con.FromPlugin <- event
}
