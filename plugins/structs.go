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
	Type    int
	Cmd     string
	Message string //todo: remove these extra fields, and just let the plugin handle based on category?
	Raw     string
	Time    time.Time
	Who     *Who
	Source  *Plugin
}

type Who struct {
	Id          string
	Name        string
	Permissions string
}
