package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/haveachin/infrared/internal/app/infrared"
)

type payloadFunc func(message string) io.Reader

type FormatterConfig struct {
	Type       string            `mapstructure:"type"`
	MessageMap map[string]string `mapstructure:"messageMap"`
}

type formatter struct {
	config           FormatterConfig
	payloadFormatter payloadFunc
}

func CreateFormatter(config FormatterConfig) *formatter {
	return &formatter{
		config:           config,
		payloadFormatter: forService(config.Type),
	}
}

func EventMessageTemplate(e EventData) map[string]string {
	return map[string]string{
		"edition":              e.Edition,
		"gatewayId":            e.GatewayID,
		"conn.network":         e.Conn.Network,
		"conn.localAddress":    e.Conn.LocalAddr,
		"conn.remoteAddress":   e.Conn.RemoteAddr,
		"conn.username":        e.Conn.Username,
		"server.serverId":      e.Server.ServerID,
		"server.serverAddress": e.Server.ServerAddr,
		"server.domains":       strings.Join(e.Server.Domains, ", "),
		"isLoginRequest":       strconv.FormatBool(*e.IsLoginRequest),
	}
}

// apply templating for available payload fields
func (f formatter) apply(e EventData, payload string) string {
	return infrared.ApplyTemplates(EventMessageTemplate(e))(payload)
}

// construct a specific endpoints payload or use the infrared default payload format
func (f formatter) payload(e EventLog) io.Reader {
	if val, ok := f.config.MessageMap[e.Type]; ok && f.payloadFormatter != nil {
		return f.payloadFormatter(f.apply(e.Data, val))
	} else {
		b, _ := json.Marshal(e.Data)
		return bytes.NewReader(b)
	}
}

// payloads for different services
func forService(name string) payloadFunc {
	switch name {
	case "discord":
		return DiscordPayload
	case "slack":
		return SlackPayload
	default:
		return nil
	}
}

func DiscordPayload(message string) io.Reader {
	payload := struct {
		Content string `json:"content"`
	}{
		Content: message,
	}
	b, _ := json.Marshal(payload)
	return bytes.NewReader(b)
}

func SlackPayload(message string) io.Reader {
	payload := struct {
		Text string `json:"text"`
	}{
		Text: message,
	}
	b, _ := json.Marshal(payload)
	return bytes.NewReader(b)
}
