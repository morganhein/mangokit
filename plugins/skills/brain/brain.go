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
	fromApp chan *plugins.Event
	toApp   chan *plugins.Event
}

var _b *brain
var config conf

func init() {
	_b = &brain{}
	plugins.RegisterPlugin("brain", plugins.Skill, _b)
	config = conf{}
}

func (b *brain) NewEvent(e plugins.Event) {
	if strings.ToLower(e.Cmd) == "quit" {
		e.Context.Say("You got it boss. Bye!")
		plugins.Core.Quit()
	}
}

func (b *brain) Setup(c *plugins.Connection) error {
	c.Events = []int{events.BOTCMD}
	return nil
}

func (b *brain) Start() error {

	return nil
}

func (b *brain) LoadConfig(location string) error {
	return nil
}
