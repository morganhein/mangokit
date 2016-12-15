package skills

import (
	"github.com/morganhein/mangokit/plugins"
	"github.com/morganhein/mangokit/events"
	"strings"
)

type conf struct {
	Token string
}

type brain struct {
	fromApp chan *plugins.Event
	toApp   chan *plugins.Event
}

var _b *brain
var config conf

func init() {
	_b = &brain{}
	plugins.RegisterSkillPlugin("brain", _b)
	config = conf{}
}

func (b *brain) NewEvent(e plugins.Event) {
	if strings.ToLower(e.Cmd) == "quit" {
		e.Context.Say("You got it boss. Bye!")
		plugins.Core.Quit()
	}
}

func (b *brain) Setup(c *plugins.Connection) ([]int, error) {
	return []int{events.BOTCMD}, nil
}

func (b *brain) LoadConfig(location string) (error) {
	return nil
}