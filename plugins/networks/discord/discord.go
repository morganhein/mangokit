package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/BurntSushi/toml"
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
)

// the plugin level server object so we can reference it when events are fired
var server *discord
// events we receive from the main program
var fromApp chan plugins.Event
// events we want to send to the main program
var toApp chan plugins.Event
// the config
var config conf

type conf struct {
	Username string
	Password string
	Token    string
}

type discord struct {
	*plugins.Plugin
	session *discordgo.Session
	me      *discordgo.User
}

func init() {
	server = &discord{}
	plugins.RegisterNetworkPlugin("discord", server)
}

func (d *discord) Setup(c *plugins.Connection) (error) {
	fromApp = c.ToPlugin
	toApp = c.FromPlugin
	// Load configuration
	return server.LoadConfig(c.Dir)
}

func (d *discord) LoadConfig(location string) (error) {
	log.Debug("Stuff")
	// todo: handle username/pass and OAuth2
	log.Debug("Loading configuration from " + location)
	if _, err := toml.DecodeFile(location, &config); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + config.Token)
	return nil
}

func (d *discord) Connect() (err error) {
	log.Debug("Connecting to Discord.")
	// Connect to discord
	d.session, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get the account information.
	d.me, err = server.session.User("@me")
	if err != nil {
		return
	}

	// Attach event handlers
	d.session.AddHandler(onConnect)
	d.session.AddHandler(onMessage)

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
	if err := d.Connect(); err != nil {
		return err
	}
	return nil
}

func (d *discord) Connected() bool {
	// todo: actually implement this
	return true
}

func (d *discord) GetContext(search interface{}) (plugins.Contexter, error) {
	return nil, nil
}

func (d *discord) getChannelByName(guild, channel string) (*discordgo.Channel, error) {
	guilds := d.session.State.Guilds
	// If we haven't yet received info on the guilds yet, populate it
	if (len(guilds) == 0) {
		// todo: figure out how to load guilds on demand, or wait until the guild info has
		// been received and then return the message
		return nil, errors.New("Guild information not yet received. Please wait and try again.")
	}
	for _, g := range guilds {
		if g.Name == guild {
			for _, ch := range g.Channels {
				if (ch.Name == channel) {
					return ch, nil
				}
			}
		}
	}
	return nil, errors.New("Could not find a guild/channel combination with that information.")
}
