package skills

import (
	"net/url"
	"fmt"
	"encoding/json"
	"time"
	"net/http"
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

func requestThought(s string) (*Thought, error) {
	url := "https://api.api.ai/v1/query?lang=en&v=20150910&sessionId=123&query=" + url.QueryEscape(s)
	fmt.Println("Requesting: " + url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer 858c143bf9a64f568bdc536b8e06697a")
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