package smalltalk

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
)

type thought struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
		Source        string `json:"source"`
		ResolvedQuery string `json:"resolvedQuery"`
		Action        string `json:"action"`
		Parameters    struct {
			Simplified string `json:"simplified"`
		} `json:"parameters"`
		Fulfillment struct {
			Speech string `json:"speech"`
		} `json:"fulfillment"`
	} `json:"result"`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
	SessionID string `json:"sessionId"`
}

type config struct {
	Token string
	api   string
}

type smalltalk struct {
	*plugins.Plugin
	conn   *plugins.Connection
	config config
}

var server *smalltalk

func init() {
	server = &smalltalk{}
	plugins.RegisterPlugin("smalltalk", plugins.Skill, server)
}

func (s *smalltalk) Setup(c *plugins.Connection) error {
	err := s.LoadConfig(c.Dir)
	c.Events = []int{events.BOTCMD}
	if err != nil {
		return err
	}
	s.conn = c
	s.config.api = "https://api.api.ai/v1/query?lang=en&v=20150910&sessionId=123&query="
	return nil
}

func (s *smalltalk) Start() error {
	log.Debug("Smalltalk started.")
	for {
		select {
		case e := <-s.conn.ToPlugin:
			if strings.HasPrefix(e.Cmd, "?") {
				msg := e.Cmd[1:]
				log.Debug("Requesting a new thought.")
				t, err := s.requestThought(s.config.api, msg)
				if err != nil {
					log.Error("Error retrieving a thought. Try thinking harder.")
					return err
				}
				e.Context.Say(t.Result.Fulfillment.Speech)
			}
		}
	}
}

func (s *smalltalk) LoadConfig(location string) error {
	log.Debug("Loading configuration from " + location)
	if _, err := toml.DecodeFile(location, &s.config); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + s.config.Token)
	return nil
}

func (s *smalltalk) requestThought(api, msg string) (*thought, error) {
	// The URL parsing is so the mockHTTP testing works
	u, err := url.Parse(api)
	u.Path = path.Join(u.Path, url.QueryEscape(msg))

	log.Debug("Requesting: " + u.String())

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+s.config.Token)
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var t thought
	err = json.NewDecoder(response.Body).Decode(&t)

	if err != nil {
		return nil, err
	}
	log.Debug("Received thought: " + t.Result.Fulfillment.Speech)
	return &t, nil
}
