package api

// Login are the requires login credentials
type Login struct {
	Username string `json:"username" validate:"required" jsonschema:"minLength=1"`
	Password string `json:"password" validate:"required" jsonschema:"minLength=1"`
}
