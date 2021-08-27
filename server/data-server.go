package server

import (
	"encoding/json"
	"github.com/ruraomsk/newServer/setup"
)

type MessageServerInfo struct {
	Status    string    `json:"status"`
	ID        int       `json:"id"`
	Interval  int       `json:"interval"`
	Channel   Channel   `json:"channel"`
	DebugInfo DebugInfo `json:"debug,omitempty"`
}
type Channel struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func newMessageServerInfo() *MessageServerInfo {
	m := new(MessageServerInfo)
	m.DebugInfo = DebugInfo{Make: setup.Set.CommServer.Debug, IP: setup.Set.CommServer.IPDebug, Port: setup.Set.CommServer.PortDebug}
	m.Channel = Channel{IP: setup.Set.CommServer.IPDevice, Port: setup.Set.CommServer.PortDevice}
	m.ID = setup.Set.CommServer.ID
	return m
}
func (m *MessageServerInfo) Adopt() []byte {
	m.Status = "Adopt"
	r, _ := json.Marshal(m)
	return r
}
func (m *MessageServerInfo) Reject() []byte {
	m.Status = "Reject"
	r, _ := json.Marshal(m)
	return r
}
