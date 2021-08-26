package server

type MessageDeviceInfo struct {
	DeviceInfo DeviceInfo `json:"deviceinfo"`
}
type DeviceInfo struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` //ETH или GPRS
	OnBoard   OnBoard   `json:"onboard,omitempty"`
	Software  SoftWare  `json:"software,omitempty"`
	Secret    Secret    `json:"secret,omitempty"`
	DebugInfo DebugInfo `json:"debug,omitempty"`
}
type OnBoard struct {
	GPRS GPRSInfo `json:"GPRS,omitempty"`
	GPS  GPSInfo  `json:"GPS,omitempty"`
	ETH  ETHInfo  `json:"ETH,omitempty"`
	USB  USBInfo  `json:"USB,omitempty"`
}
type SoftWare struct {
	Version string `json:"version"`
}
type Secret struct {
	Key string `json:"key"`
}
type DebugInfo struct {
	Make bool   `json:"make"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type GPRSInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type GPSInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type ETHInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
type USBInfo struct {
	Present bool   `json:"present"`
	Model   string `json:"model"`
	Status  string `json:"status"`
}
