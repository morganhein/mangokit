package mangokit

import (
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
	_ "github.com/morganhein/mangokit/plugins/networks/discord"
)


/*
TODO: implement permissions system
TODO: logging for all actions, to a log and to the screen
TODO: implement plugging into twitch.tv, discord, etc
*/

var b *Brain

func init() {
	b = &Brain{}
}

func Start() () {
	log.Debug("Bootstrap has begun.")

	// Load the basic configuration


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
	log.Debug("Bootstrapping finished.")
	b.Loop()
}