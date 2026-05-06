package core_http_response

type ErrorsResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
