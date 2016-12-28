package smalltalk

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/morganhein/mangokit/events"
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
	conf config
}

var log = plugins.GetLogger()
var server *smalltalk

func init() {
	server = &smalltalk{
		Plugin: plugins.NewPlugin("smalltalk", plugins.Skill, []int{events.BOTCMD}),
	}
	plugins.RegisterPlugin(server)
}

func (s *smalltalk) Start() error {
	log.Debug("Smalltalk started.")
	err := s.LoadConfig(s.Dir())
	if err != nil {
		log.Error("Unable to load configuration for " + s.Name())
	}
	s.conf.api = "https://api.api.ai/v1/query?lang=en&v=20150910&sessionId="

	for {
		select {
		case e := <-s.ToPlugin():
			if strings.HasPrefix(e.Cmd, "?") {
				msg := e.Cmd[1:]
				log.Debug("Requesting a new thought.")
				t, err := s.requestThought(s.conf.api, e.Who.Id, msg)
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
	if _, err := toml.DecodeFile(location, &s.conf); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + s.conf.Token)
	return nil
}

func (s *smalltalk) requestThought(api, session, msg string) (*thought, error) {
	// The URL parsing is so the mockHTTP testing works
	path := api + session + "&query=" + url.QueryEscape(msg)

	log.Debug("Requesting: " + path)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", "Bearer "+s.conf.Token)
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
