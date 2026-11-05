package api

// RTMPChannel represents details about a currently connected RTMP publisher
type RTMPChannel struct {
	Name string `json:"name" jsonschema:"minLength=1"`
}
