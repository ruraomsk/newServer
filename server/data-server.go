package server

type MessageServerInfo struct {
	Status    string    `json:"status"`
	ID        int       `json:"id"`
	Timeout   int       `json:"timeout"`
	Channel   Channel   `json:"channel"`
	DebugInfo DebugInfo `json:"debug,omitempty"`
}
type Channel struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
