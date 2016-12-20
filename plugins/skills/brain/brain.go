package skills

import (
	"strings"

	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/plugins"
)

type conf struct {
	Token string
}

type brain struct {
	*plugins.Plugin
	conn *plugins.Connection
}

var _b *brain
var config conf

func init() {
	_b = &brain{}
	plugins.RegisterPlugin("brain", plugins.Skill, _b)
	config = conf{}
}

func (b *brain) Setup(c *plugins.Connection) error {
	c.Events = []int{events.BOTCMD}
	b.conn = c
	return nil
}

func (b *brain) Start() error {
	for {
		select {
		case e := <-b.conn.ToPlugin:
			if strings.HasPrefix(e.Cmd, "!") {
				msg := e.Cmd[1:]
				switch msg {
				case "quit":
					e.Context.Say("Ok fine. Bye!")
					plugins.Shutdown("Requested shutdown by " + e.Who.Name)
					return nil
				}
			}
		}
	}
	return nil
}

func (b *brain) LoadConfig(location string) error {
	return nil
}
