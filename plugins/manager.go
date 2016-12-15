package plugins

import (
	"os"
	"github.com/morganhein/mangokit/log"
	"path"
	"strings"
	"github.com/morganhein/mangokit/events"
)

var Core CoreFramework
var NetworkPlugins = make(map[*Connection]NetworkPlugineers)
var SkillPlugins = make(map[*Connection]SkillPlugineers)

func RegisterNetworkPlugin(name string, plugin NetworkPlugineers) {
	c := registerNewPlugin(name, Network)
	if c != nil {
		NetworkPlugins[c] = plugin
	}
}

func RegisterSkillPlugin(name string, plugin SkillPlugineers) {
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
		dir = path.Join(dir, "plugins", "networks", name, "config.toml")
	case Skill:
		dir = path.Join(dir, "plugins", "skills", name, "config.toml")
	}
	log.Debug("Creating plugin from dir: " + dir)

	c := &Connection{
		Name: name,
		ToPlugin: toPlugin,
		FromPlugin:fromPlugin,
		Type: t,
		Dir: dir,
	}
	return c
}

func PopulateCmd(e *Event) (error) {
	// Is this a message?
	if e.Type != events.PRIVATEMESSAGE && e.Type != events.PUBLICMESSAGE && e.Type != events.MESSAGE {
		return nil
	}
	if !strings.HasPrefix(e.Raw, ".") {
		e.Message = e.Raw
		return nil
	}
	if !strings.Contains(e.Raw, ":") {
		e.Message = e.Raw
		return nil
	}
	split := strings.SplitN(e.Raw, ":", 2)

	// sanity checking here, not sure it's required
	if len(split) != 2 {
		e.Message = e.Raw
		return nil
	}
	e.Cmd = split[0][1:]
	e.Message = split[1]
	e.Type = events.BOTCMD
	log.Debug("Found a new command: " + e.Cmd)
	return nil
}