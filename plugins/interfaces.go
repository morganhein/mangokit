package plugins

type CoreFramework interface {
	Leave(Contexter)
	Quit()
}

// Contexter is the interface by which all events should have a reference to, so they can respond back
type Contexter interface {
	// is this context currently messageable (offline, or the base server context which can't rx messages)
	Messageable() bool
	Say(message string) error
	Name() string
	Who() []Contexter
	// Send a file
	// File() (sometimes a pseudonym for Attach() ?

	// Attach a file
	// Attach()

	// Stream something (video, audio, whatever)
	// StreamAudio()
	// StreamVideo()

	// Emoticon support
	// Emoticon(name string)

	// Send a raw command
	// Command(cmd string) error
}

// Plugineers is the interface for all Plugins
type Plugineer interface {
	// Start/Connect the network service/skill (will be run via go Start())
	Start() error
	// Get a context for a channel, user, board, group, guild, etc
	// GetContext(interface{}) (Contexter, error) //todo: later if desired
	// Disconnect from the network service
	Disconnect() error
	// Status of the connection
	Connected() bool
	// Disconnect if still connected, and reconnect
	Reconnect() error
	// Force load a configuration file. This normally should only be called by Setup()
	LoadConfig(location string) error
	// Shutdown all running processes gracefully
	Shutdown()
	// Name returns the name of this plugin
	Name() string
	// Dir returns the directory this plugin exists under
	Dir() string
	// Events returns the events this plugin wants to listen to
	Events() []int
	// Category returns the type of this plugin (Network or Skill)
	Category() int
	// ToPlugin returns the channel to send events to this plugin
	ToPlugin() chan Event
	// FromPlugin returns the channel for receiving events from the plugin
	FromPlugin() chan Event
}

// Logger is the basic interface which all logging should be sent to
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Notice(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Critical(...interface{})
}
