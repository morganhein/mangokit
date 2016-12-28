package mangokit

import (
	"github.com/morganhein/mangokit/plugins"
	_ "github.com/morganhein/mangokit/plugins/networks/discord" // all blank imports register the plugin
	_ "github.com/morganhein/mangokit/plugins/skills/brain"
	_ "github.com/morganhein/mangokit/plugins/skills/smalltalk"
)

/*
TODO: implement permissions system
TODO: logging for all actions, to a log and to the screen
TODO: implement plugging into twitch.tv, discord, etc
*/

var log plugins.Logger

func Start() {
	log = plugins.GetLogger()

	log.Debug("Bootstrap has begun.")

	// Load the basic configuration
	// todo: loading here

	// iterate each network plugin
	for _, p := range plugins.NetworkPlugins {
		log.Debug("Setting up plugin " + p.Name())
		go p.Start()
	}
	for _, p := range plugins.SkillPlugins {
		log.Debug("Setting up plugin " + p.Name())
		Core.AddEventTriggers(p)
		go p.Start()
	}
	log.Debug("Bootstrapping finished.")
	Core.Loop()
}
