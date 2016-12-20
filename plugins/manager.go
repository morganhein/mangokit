package plugins

import (
	"os"
	"strings"

	"path/filepath"

	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/log"
)

var Core CoreFramework
var NetworkPlugins = make(map[*Connection]Plugineers)
var SkillPlugins = make(map[*Connection]Plugineers)

func RegisterPlugin(name string, pluginType int, plugin Plugineers) {
	switch pluginType {
	case 0:
		registerNetworkPlugin(name, plugin)
	case 1:
		registerSkillPlugin(name, plugin)
	}
}

func registerNetworkPlugin(name string, plugin Plugineers) {
	c := registerNewPlugin(name, Network)

	if c != nil {
		NetworkPlugins[c] = plugin
		c.Plugin = plugin
	}
}

func registerSkillPlugin(name string, plugin Plugineers) {
	c := registerNewPlugin(name, Skill)

	if c != nil {
		SkillPlugins[c] = plugin
	}
}

func registerNewPlugin(name string, t int) *Connection {
	fromPlugin := make(chan Event, 10)
	toPlugin := make(chan Event, 10)

	// todo: enable this later, for actual builds
	//dir, err := osext.ExecutableFolder();
	//Brain.log.Debug(dir)
	// todo: for now, we'll use this in our build environment
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Unable to find the current working directory.")
		return nil
	}

	switch t {
	case Network:
		dir = filepath.Join(dir, "plugins", "networks", name, "config.toml")
	case Skill:
		dir = filepath.Join(dir, "plugins", "skills", name, "config.toml")
	}
	log.Debug("Creating plugin from dir: " + dir)

	c := &Connection{
		Name:       name,
		ToPlugin:   toPlugin,
		FromPlugin: fromPlugin,
		Type:       t,
		Dir:        dir,
	}
	return c
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
