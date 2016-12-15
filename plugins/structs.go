package plugins

import (
	"time"
)

const (
	Network = iota
	Skill
)

type Event struct {
	Context    Contexter
	// Type is event Type (verb), like Connected, Messaged, Quit, etc.
	Type       int
	Cmd        string
	Message    string
	Raw        string
	Time       time.Time
	Who        *Who
	Connection *Connection
}

type Who struct {
	Id          string
	Name        string
	Permissions string
}

type Connection struct {
	// Name of the plugin. The plugin MUST reside under /plugins/networks/<name>/ or /plugins/skills/<name>/
	Name       string
	// Send to plugin
	ToPlugin   chan Event
	// Receive from plugin
	FromPlugin chan Event
	// Directory the plugin exists in
	Dir        string
	// The Type of the plugin
	Type       int
}
