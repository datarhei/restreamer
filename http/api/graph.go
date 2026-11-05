package api

type GraphQuery struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type GraphResponse struct {
	Data   interface{}   `json:"data"`
	Errors []interface{} `json:"errors"`
}
