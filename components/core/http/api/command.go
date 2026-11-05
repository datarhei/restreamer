package api

// Command is a command to send to a process
type Command struct {
	Command string `json:"command" validate:"required" enums:"start,stop,restart,reload" jsonschema:"enum=start,enum=stop,enum=restart,enum=reload"`
}
