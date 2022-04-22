package callback

import "fmt"

const (
	EventTypeError          string = "Error"
	EventTypePlayerJoin     string = "PlayerJoin"
	EventTypePlayerLeave    string = "PlayerLeave"
	EventTypeContainerStart string = "ContainerStart"
	EventTypeContainerStop  string = "ContainerStop"
)

type Event interface {
	EventType() string
	EventMsg() string
}

type ErrorEvent struct {
	Error    string `json:"error"`
	ProxyUID string `json:"proxyUid"`
}

func (event ErrorEvent) EventType() string {
	return EventTypeError
}

func (event ErrorEvent) EventMsg() string {
	return fmt.Sprintf("%s @%s", event.Error, event.ProxyUID)
}

type PlayerJoinEvent struct {
	Username      string `json:"username"`
	RemoteAddress string `json:"remoteAddress"`
	TargetAddress string `json:"targetAddress"`
	ProxyUID      string `json:"proxyUid"`
}

func (event PlayerJoinEvent) EventType() string {
	return EventTypePlayerJoin
}

func (event PlayerJoinEvent) EventMsg() string {
	return fmt.Sprintf("%s has connected to %s", event.Username, event.TargetAddress)
}

type PlayerLeaveEvent struct {
	Username      string `json:"username"`
	RemoteAddress string `json:"remoteAddress"`
	TargetAddress string `json:"targetAddress"`
	ProxyUID      string `json:"proxyUid"`
}

func (event PlayerLeaveEvent) EventType() string {
	return EventTypePlayerLeave
}

func (event PlayerLeaveEvent) EventMsg() string {
	return fmt.Sprintf("%s has left %s", event.Username, event.TargetAddress)
}

type ContainerStartEvent struct {
	ProxyUID string `json:"proxyUid"`
}

func (event ContainerStartEvent) EventType() string {
	return EventTypeContainerStart
}

func (event ContainerStartEvent) EventMsg() string {
	return fmt.Sprintf("infrared has started %s", event.ProxyUID)
}

type ContainerStopEvent struct {
	ProxyUID string `json:"proxyUid"`
}

func (event ContainerStopEvent) EventType() string {
	return EventTypeContainerStop
}

func (event ContainerStopEvent) EventMsg() string {
	return fmt.Sprintf("infrared has stopped %s", event.ProxyUID)
}
