package plugins

// Plugin is the base type for all plugins
type Plugin struct{}

// Config holds all plugin configuration information.
type config struct {
	Dir string
}

var this Plugin
var conf config

//func init() {
//	this := &Plugin{}
//	RegisterPlugin("plugin", Network, this)
//	// conf := &config{}
//}

// Setup stores a reference to the connection after filling the configuration information.
func (p *Plugin) Setup(c *Connection) error {
	return p.LoadConfig(c.Dir)
}

// Start should begin the main loop of the plugin
func (p *Plugin) Start() error {
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
	Log.Debug("Loading configuration from " + location)
	return nil
}

func (p *Plugin) Shutdown() {

}
