package plugins

type CoreFramework interface {
	Leave(Contexter)
	Quit()
}

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

type Plugineers interface {
	// Setup plugin with bootstrap info, should also load configuration
	Setup(*Connection) error
	// Start/Connect the network service/skill (will be run via go Start())
	Start() error
	// Get a context for a channel, user, board, group, guild, etc
	GetContext(interface{}) (Contexter, error)
	// Disconnect from the network service
	Disconnect() error
	// Status of the connection
	Connected() bool
	// Disconnect if still connected, and reconnect
	Reconnect() error
	// Force load a configuration file. This normally should only be called by Setup()
	LoadConfig(location string) error
}

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Notice(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Critical(...interface{})
}
