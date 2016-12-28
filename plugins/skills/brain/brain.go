package skills

import (
	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/plugins"
)

type config struct {
	Token string
}

type brain struct {
	*plugins.Plugin
	conf config
}

var log = plugins.GetLogger()
var _b *brain

func init() {
	_b = &brain{
		Plugin: plugins.NewPlugin("brain", plugins.Skill, []int{events.BOTCMD}),
	}
	plugins.RegisterPlugin(_b)
}

func (b *brain) Start() error {
	log.Debug("Brain started.")

	for {
		select {
		case e := <-b.ToPlugin():
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

func (b *brain) LoadConfig(location string) error {
	return nil
}
