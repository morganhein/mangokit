package discord

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/plugins"
)

type discord struct {
	*plugins.Plugin
	// events we want to send/receive to/from the main program
	// saved globally due to the way DiscordGo passes events
	session *discordgo.Session
	me      *discordgo.User
	conf    config
}

type config struct {
	Username string
	Password string
	Token    string
	Owner    string
}

var disc *discord
var log = plugins.GetLogger()

func init() {
	disc = &discord{
		Plugin: plugins.NewPlugin("discord", plugins.Network, []int{events.MESSAGE}),
	}
	plugins.RegisterPlugin(disc)
}

func (d *discord) LoadConfig(location string) error {
	log.Debug("Loading configuration from " + location)
	if _, err := toml.DecodeFile(location, &d.conf); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + d.conf.Token)
	return nil
}

func (d *discord) Start() (err error) {
	err = d.LoadConfig(d.Dir())
	if err != nil {
		log.Error("Unable to load configuration for " + d.Name())
	}
	log.Debug("Connecting to Discord.")
	// Connect to discord
	d.session, err = discordgo.New("Bot " + d.conf.Token)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get the account information.
	d.me, err = d.session.User("@me")
	//server.session.User("@me")
	if err != nil {
		return
	}

	// Attach event handlers
	d.session.AddHandler(d.onConnect)
	d.session.AddHandler(d.onMessage)
	d.session.AddHandler(d.onMessageDelete)
	d.session.AddHandler(d.onMessageUpdate)
	d.session.AddHandler(d.onGuildMemberAdd)
	d.session.AddHandler(d.onGuildMemberExit)

	log.Debug("Discord started.")
	// Open the websocket connection.
	err = d.session.Open()

	return
}

func (d *discord) Disconnect() error {
	return d.session.Close()
}

func (d *discord) Reconnect() error {
	if err := d.Disconnect(); err != nil {
		return err
	}
	if err := d.connect(); err != nil {
		return err
	}
	return nil
}

func (d *discord) connect() error {
	return d.session.Open()
}

func (d *discord) Connected() bool {
	// todo: actually implement this
	return true
}

func (d *discord) getChannelByName(guild, channel string) (*discordgo.Channel, error) {
	guilds := d.session.State.Guilds
	// If we haven't yet received info on the guilds yet, populate it
	if len(guilds) == 0 {
		// todo: figure out how to load guilds on demand, or wait until the guild info has
		// been received and then return the message
		return nil, errors.New("Guild information not yet received. Please wait and try again.")
	}
	for _, g := range guilds {
		if g.Name == guild {
			for _, ch := range g.Channels {
				if ch.Name == channel {
					return ch, nil
				}
			}
		}
	}
	return nil, errors.New("Could not find a guild/channel combination with that information.")
}

func (d *discord) Shutdown() {
	d.session.Close()
}
