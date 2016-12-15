package mangokit

import (
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
	_ "github.com/morganhein/mangokit/plugins/networks/discord"
	_ "github.com/morganhein/mangokit/plugins/skills/smalltalk"
)


/*
TODO: implement permissions system
TODO: logging for all actions, to a log and to the screen
TODO: implement plugging into twitch.tv, discord, etc
*/

func Start() () {
	log.Debug("Bootstrap has begun.")

	// Load the basic configuration
	// todo: loading here

	// iterate each network plugin
	for c, p := range plugins.NetworkPlugins {
		log.Debug("Setting up plugin " + c.Name)
		// Setup the plugin
		if err := p.Setup(c); err != nil {
			log.Critical(err.Error())
			continue
		}
		// Have the p connect
		if err := p.Connect(); err != nil {
			log.Critical(err.Error())
			continue
		}
	}
	for c, p := range plugins.SkillPlugins {
		log.Debug("Setting up plugin " + c.Name)
		// Setup the plugin
		es, err := p.Setup(c)
		if err != nil {
			log.Critical(err.Error())
			continue
		}
		Brain.AddEventTriggers(es, p)
	}
	log.Debug("Bootstrapping finished.")
	Brain.Loop()
}