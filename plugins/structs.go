package plugins

import (
	"time"
)

const (
	Network = iota
	Skill
)

type Event struct {
	Context Contexter
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
	Name string
	// Reference to the plugin implementation
	Plugin Plugineers
	// ToPlugin receives events from the main handler
	ToPlugin chan Event
	// FromPlugin sends events to the main handler
	FromPlugin chan Event
	// Dir is the directory the plugin exists in
	Dir string
	// Type of the plugin (0:Network,1:Skill)
	Type int
	// Events is the list of events the plugin wants to receive.
	Events []int
}
