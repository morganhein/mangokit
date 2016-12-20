package smalltalk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morganhein/mangokit/plugins"
	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	ret string
}

func (c *mockContext) Say(msg string) error {
	c.ret = msg
	return nil
}

func (c *mockContext) Who() []plugins.Contexter {
	c.ret = "who"
	return nil
}

func (c *mockContext) Messageable() bool {
	return true
}

func (c *mockContext) Name() string {
	return "name"
}

func TestThoughtRequest(t *testing.T) {
	s := smalltalk{}

	th := thought{}
	th.Result.Fulfillment.Speech = "This is a test!"
	mockHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(th)
	}))
	defer mockHTTP.Close()

	response, err := s.requestThought(mockHTTP.URL, "This is a dummy request.")
	assert.Nil(t, err)
	assert.Equal(t, response.Result.Fulfillment.Speech, "This is a test!")
}

func TestSendingMessage(t *testing.T) {
	th := thought{}
	th.Result.Fulfillment.Speech = "This is a test!"
	mockHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(th)
	}))
	defer mockHTTP.Close()

	s := smalltalk{}
	s.config = config{
		api: mockHTTP.URL,
	}

	e := plugins.Event{}
	e.Context = &mockContext{}

	go s.Start()
}
