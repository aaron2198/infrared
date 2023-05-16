package webhook

import (
	"errors"
	"net/http"
	"time"
)

var ErrEventTypeNotAllowed = errors.New("event topic not allowed")

// HTTPClient represents an interface for the Webhook to send events with.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type EventData struct {
	Edition   string `json:"edition"`
	GatewayID string `json:"gatewayId"`
	Conn      struct {
		Network    string `json:"network"`
		LocalAddr  string `json:"localAddress"`
		RemoteAddr string `json:"remoteAddress"`
		Username   string `json:"username,omitempty"`
	} `json:"client"`
	Server struct {
		ServerID   string   `json:"serverId,omitempty"`
		ServerAddr string   `json:"serverAddress,omitempty"`
		Domains    []string `json:"domains,omitempty"`
	} `json:"server"`
	IsLoginRequest *bool `json:"isLoginRequest,omitempty"`
}

// EventLog is the struct that will be send to the Webhook.URL
type EventLog struct {
	Type       string    `json:"type"`
	Topics     []string  `json:"topics"`
	OccurredAt time.Time `json:"occurredAt"`
	Data       EventData `json:"data"`
}

// Webhook can send a Event via POST Request to a specified URL.
// There are two ways to use a Webhook. You can directly call
// DispatchEvent or Serve to attach a channel to the Webhook.
type Webhook struct {
	ID             string
	HTTPClient     HTTPClient
	URL            string
	AllowedTopics  []string
	AllowedServers []string
	Formatter      *formatter
}

// hasEvent checks if Webhook.EventTypes contain the given event's type.
func (webhook Webhook) hasEvent(e EventLog) bool {
	inServerList := false
	if len(webhook.AllowedServers) != 0 {
		for _, as := range webhook.AllowedServers {
			if as == e.Data.Server.ServerID {
				inServerList = true
			}
		}
	}

	if inServerList {
		for _, at := range webhook.AllowedTopics {
			for _, et := range e.Topics {
				if at == et {
					return true
				}
			}
		}
	}
	return false
}

// DispatchEvent wraps the given Event in an EventLog and marshals it into JSON
// before sending it in a POST Request to the Webhook.URL.
func (webhook Webhook) DispatchEvent(e EventLog) error {
	if !webhook.hasEvent(e) {
		return ErrEventTypeNotAllowed
	}

	request, err := http.NewRequest(http.MethodPost, webhook.URL, webhook.Formatter.payload(e))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := webhook.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	// We don't care about the client's response, but we should still close the client's body if it exists.
	// If not closed the underlying connection cannot be reused for further requests.
	// See https://pkg.go.dev/net/http#Client.Do for more details.
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	return nil
}
