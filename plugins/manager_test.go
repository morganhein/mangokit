package plugins

import (
	"github.com/morganhein/mangokit/events"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPopulateCmd(t *testing.T) {
	event := Event{
		Type: events.MESSAGE,
		Raw: ".yt:attribute",
	}
	err := PopulateCmd(&event)
	assert.Equal(t, err, nil, "Should not return an error.")
	assert.Equal(t, "yt", event.Cmd, "Command should be yt.")
	assert.Equal(t, "attribute", event.Message, "Attribute should be attribute.")
	assert.Equal(t, events.BOTCMD, event.Type, "Event type should be a botmessage.")

}