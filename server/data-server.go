package server

import (
	"github.com/ruraomsk/newServer/setup"
)

type MessageServerInfo struct {
	Status   string  `json:"status"`
	ID       int     `json:"id"`
	Interval int     `json:"interval"`
	Channel  Channel `json:"channel"`
}
type Channel struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type Confirm struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
type MessageServer struct {
	Connect MessageServerInfo `json:"connect,omitempty"`
}
type MessageDevice struct {
	Connect DeviceInfo `json:"connect,omitempty"`
	Confirm Confirm    `json:"confirm,omitempty"`
}

func newMessageServerInfo() *MessageServerInfo {
	m := new(MessageServerInfo)
	m.Channel = Channel{IP: setup.Set.CommServer.IPDevice, Port: setup.Set.CommServer.PortDevice}
	m.ID = setup.Set.CommServer.ID
	return m
}
func (m *MessageServerInfo) Adopt() {
	m.Status = "Adopt"
}
func (m *MessageServerInfo) Reject() {
	m.Status = "Reject"
}
