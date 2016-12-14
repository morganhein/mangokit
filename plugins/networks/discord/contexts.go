package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"morganhein/mangobot/plugins"
)

type guild struct {
	s *discordgo.Session
}

func (g *guild) Say(message string) error {
	return errors.New("Cannot message an entire Guild.")
}

func (g *guild) Messageable() bool {
	return false
}

func (g *guild) Name() string {
	return "Guild Name Here"
}

// Guild is synonymous with "server" in the discord context
func (g *guild) Who() []plugins.Contexter {
	chans := make([]plugins.Contexter, 10)
	//for i, v := range chans {
	//	chans[i] = plugins.Contexter(&v)
	//}
	return chans
}

type channel struct {
	s *discordgo.Session
	c *discordgo.Channel
}

func (c *channel) Say(message string) (err error) {
	_, err = c.s.ChannelMessageSend(c.c.ID, message)
	return
}

func (c *channel) Messageable() bool {
	//todo detect here if the current channel can receive messages
	return true
}

func (c *channel) Name() string {
	return c.c.Name
}

func (c *channel) Who() []plugins.Contexter {
	members := make([]plugins.Contexter, 10)
	//for i, v := range members {
	//	members[i] = plugins.Contexter(&v)
	//}
	return members
}

type member struct {
	s *discordgo.Session
	m *discordgo.Member
}

func (m *member) Say(message string) (err error) {
	// list all user channels
	chans, err := m.s.UserChannels()
	cid := ""
	for _, c := range chans {
		// if the private channel for this member already exists use it
		if c.Recipient != nil &&
			c.Recipient.ID == m.m.User.ID &&
			c.IsPrivate {
			cid = c.ID
		}
	}
	// otherwise create a new one
	if cid == "" {
		ch, err := m.s.UserChannelCreate(m.m.User.ID)
		if err != nil {
			return err
		}
		cid = ch.ID
	}
	// and use it
	_, err = m.s.ChannelMessageSend(cid, message)
	return err
}

func (m *member) Messageable() bool {
	return true
}

func (m *member) Who() []plugins.Contexter {
	members := make([]plugins.Contexter, 1)
	//m[0] = plugins.Contexter(&m[0])
	return members
}
