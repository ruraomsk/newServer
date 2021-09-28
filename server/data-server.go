package server

import (
	"github.com/ruraomsk/ag-server/pudge"
)

type DeviceInfo struct {
	ID    int    `json:"id"`
	Type  string `json:"type"` //ETH или GPRS
	Ready bool   `json:"ready"`
}

type Confirm struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
type MessageDevice struct {
	Connect DeviceInfo   `json:"connect,omitempty"`
	Confirm Confirm      `json:"confirm,omitempty"`
	Status  DeviceStatus `json:"status,omitempty"`
}
type DeviceStatus struct {
	Ready    bool                  `json:"ready"`
	TechMode int                   `json:"tech"`
	Base     bool                  `json:"base"`
	PK       int                   `json:"pk"`
	CK       int                   `json:"ck"`
	NK       int                   `json:"nk"`
	TMax     int                   `json:"tmax"`
	ComDU    pudge.StatusCommandDU `json:"comdu"`
	DK       pudge.DK              `json:"dk"`
}

func giveMeStatus() string {
	return "give_me_status"
}
