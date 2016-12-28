package plugins

import (
	"os"
	"strings"

	"github.com/morganhein/mangokit/events"
)

var Core CoreFramework
var NetworkPlugins = make([]Plugineer, 0)
var SkillPlugins = make([]Plugineer, 0)

func RegisterPlugin(p Plugineer) {
	log.Debug("Registering plugin: " + p.Name())
	switch p.Category() {
	case Network:
		NetworkPlugins = append(NetworkPlugins, p)
	case Skill:
		SkillPlugins = append(SkillPlugins, p)
	}
}

func PopulateCmd(e *Event) error {
	// Is this a message?
	if e.Type != events.PRIVATEMESSAGE && e.Type != events.PUBLICMESSAGE && e.Type != events.MESSAGE {
		return nil
	}
	if !strings.HasPrefix(e.Raw, ".") {
		e.Message = e.Raw
		return nil
	}
	e.Cmd = e.Raw[1:]
	e.Type = events.BOTCMD
	log.Debug("Found a new command: " + e.Cmd)
	return nil
}

func Shutdown(reason string) {
	for _, p := range SkillPlugins {
		p.Shutdown()
	}
	for _, p := range NetworkPlugins {
		p.Shutdown()
	}
	//todo graceful shutdown
	os.Exit(0)
}
