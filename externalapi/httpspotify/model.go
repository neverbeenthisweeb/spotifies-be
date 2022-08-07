package httpspotify

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
