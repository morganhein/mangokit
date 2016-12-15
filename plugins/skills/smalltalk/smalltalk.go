package skills

import (
	"net/url"
	"fmt"
	"encoding/json"
	"time"
	"net/http"
	"github.com/morganhein/mangokit/plugins"
	"github.com/BurntSushi/toml"
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/events"
)

type Thought struct {
	ID        string `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
			  Source        string `json:"source"`
			  ResolvedQuery string `json:"resolvedQuery"`
			  Action        string `json:"action"`
			  Parameters    struct {
						Simplified string `json:"simplified"`
					} `json:"parameters"`
			  Fulfillment   struct {
						Speech string `json:"speech"`
					} `json:"fulfillment"`
		  } `json:"result"`
	Status    struct {
			  Code      int `json:"code"`
			  ErrorType string `json:"errorType"`
		  } `json:"status"`
	SessionID string `json:"sessionId"`
}

type conf struct {
	Token string
}

type smalltalk struct {
	fromApp chan *plugins.Event
	toApp   chan *plugins.Event
}

var server *smalltalk
var config conf

func init() {
	server = &smalltalk{}
	plugins.RegisterSkillPlugin("smalltalk", server)
	config = conf{}
}

func (s *smalltalk) NewEvent(e plugins.Event) {
	if t, err := s.requestThought(e.Message); err == nil {
		e.Context.Say(t.Result.Fulfillment.Speech)
	}
}

func (s *smalltalk) Setup(c *plugins.Connection) ([]int, error) {
	err := s.LoadConfig(c.Dir)
	if err != nil {
		return []int{}, err
	}
	return []int{events.BOTCMD}, nil
}

func (s *smalltalk) LoadConfig(location string) (error) {
	log.Debug("Loading configuration from " + location)
	if _, err := toml.DecodeFile(location, &config); err != nil {
		log.Error("Could not load configuration file: " + err.Error())
		return err
	}
	log.Debug("Loaded configuration file with Token: " + config.Token)
	return nil
}

func (s *smalltalk) requestThought(str string) (*Thought, error) {
	url := "https://api.api.ai/v1/query?lang=en&v=20150910&sessionId=123&query=" + url.QueryEscape(str)
	fmt.Println("Requesting: " + url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + config.Token)
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var thought Thought
	err = json.NewDecoder(response.Body).Decode(&thought)

	if err != nil {
		return nil, err
	}
	fmt.Println("Received thought: " + thought.Result.Fulfillment.Speech)
	return &thought, nil
}