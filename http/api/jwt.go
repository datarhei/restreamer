package api

// JWT is the JWT token and its expiry date
type JWT struct {
	AccessToken  string `json:"access_token" jsonschema:"minLength=1"`
	RefreshToken string `json:"refresh_token" jsonschema:"minLength=1"`
}

type JWTRefresh struct {
	AccessToken string `json:"access_token" jsonschema:"minLength=1"`
}
