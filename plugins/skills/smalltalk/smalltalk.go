package skills

import (
	"encoding/json"
	"net/http"
	"net/url"
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

type conf struct {
	Token string
}

type smalltalk struct {
	*plugins.Plugin
	conn *plugins.Connection
}

var server *smalltalk
var config conf

func init() {
	server = &smalltalk{}
	plugins.RegisterPlugin("smalltalk", plugins.Skill, server)
	config = conf{}
}

func (s *smalltalk) NewEvent(e plugins.Event) {
	if t, err := s.requestThought(e.Message); err == nil {
		e.Context.Say(t.Result.Fulfillment.Speech)
	}
}

func (s *smalltalk) Setup(c *plugins.Connection) error {
	err := s.LoadConfig(c.Dir)
	c.Events = []int{events.BOTCMD}
	if err != nil {
		return err
	}
	s.conn = c
	return nil
}

func (s *smalltalk) Start() error {
	log.Debug("Smalltalk started.")
	for {
		select {
		case e := <-s.conn.ToPlugin:
			log.Debug("Requesting a new thought.")
			t, err := s.requestThought(e.Cmd)
			if err != nil {
				log.Error("Error retrieving a thought. Try thinking harder.")
			}
			e.Context.Say(t.Result.Fulfillment.Speech)
		}
	}
	log.Debug("Smalltalk exiting.")
	return nil
}

func (s *smalltalk) LoadConfig(location string) error {
	log.Debug("Loading configuration from " + location)
	if _, err := toml.DecodeFile(location, &config); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + config.Token)
	return nil
}

func (s *smalltalk) requestThought(str string) (*thought, error) {
	url := "https://api.api.ai/v1/query?lang=en&v=20150910&sessionId=123&query=" + url.QueryEscape(str)
	log.Debug("Requesting: " + url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.Token)
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
