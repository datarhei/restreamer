package api

type WidgetProcess struct {
	CurrentSessions uint64 `json:"current_sessions" format:"uint64"`
	TotalSessions   uint64 `json:"total_sessions" format:"uint64"`
	Uptime          int64  `json:"uptime"`
}
