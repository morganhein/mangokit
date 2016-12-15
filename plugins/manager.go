package plugins

import (
	"os"
	"github.com/morganhein/mangokit/log"
	"path"
)

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