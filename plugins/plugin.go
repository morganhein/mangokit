package plugins

import (
	"os"
	"path/filepath"
)

// Plugin is the base type for all plugins
type Plugin struct {
	// Name of the plugin. The plugin MUST reside under /plugins/networks/<name>/ or /plugins/skills/<name>/
	name string
	// ToPlugin receives events from the main handler
	toPlugin chan Event
	// FromPlugin sends events to the main handler
	fromPlugin chan Event
	// control is the channel to send Control Messages (exit/restart/etc) to the plugin
	control chan int
	// Dir is the directory the plugin exists in
	dir string
	// Category of the plugin (0:Network,1:Skill)
	category int
	// Events is the list of events the plugin wants to receive.
	events []int
}

var this Plugin

func NewPlugin(name string, category int, events []int) *Plugin {
	// todo: enable this later, for actual builds
	//dir, err := osext.ExecutableFolder();
	//Brain.log.Debug(dir)
	// todo: for now, we'll use this in our build environment
	dir, err := os.Getwd()

	if err != nil {
		log.Error("Unable to get current working directory to load plugin.")
		return nil
	}

	switch category {

	case Network:
		dir = filepath.Join(dir, "plugins", "networks", name, "config.toml")
	case Skill:
		dir = filepath.Join(dir, "plugins", "skills", name, "config.toml")
	}

	p := &Plugin{
		category:   category,
		dir:        dir,
		events:     events,
		name:       name,
		fromPlugin: make(chan Event, 10),
		toPlugin:   make(chan Event, 10),
	}
	return p
}

// Start should begin the main loop of the plugin
func (p *Plugin) Start() error {
	p.fromPlugin = make(chan Event, 10)
	p.toPlugin = make(chan Event, 10)
	return nil
}

// Disconnect should close all network sessions
func (p *Plugin) Disconnect() error {
	return nil
}

// Connected returns true if the network is connected, false otherwise
func (p *Plugin) Connected() bool {
	return false
}

// Reconnect should disconnect and reconnect again.
func (p *Plugin) Reconnect() error {
	return nil
}

// LoadConfig forcefully loads a config from the passed location.
func (p *Plugin) LoadConfig(location string) error {
	p.dir = location
	log.Debug("Loading configuration from " + location)
	return nil
}

func (p *Plugin) Shutdown() {

}

func (p *Plugin) Name() string {
	return p.name
}

func (p *Plugin) Dir() string {
	return p.dir
}

func (p *Plugin) Events() []int {
	return p.events
}

func (p *Plugin) Category() int {
	return p.category
}

func (p *Plugin) ToPlugin() chan Event {
	return p.toPlugin
}

func (p *Plugin) FromPlugin() chan Event {
	return p.fromPlugin
}
