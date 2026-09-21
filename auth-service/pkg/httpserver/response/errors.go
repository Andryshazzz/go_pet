package httpresponse

// ErrorsResponse represents a standardized error response body.
//
// Example JSON:
//
//	{
//	    "error": "invalid FullName length 1 (must be 2-100): invalid argument",
//	    "message": "failed to create user"
//	}
type ErrorsResponse struct {
	// Error contains the technical error details for debugging.
	Error string `json:"error"`
	
	// Message contains the human-readable error description.
	Message string `json:"message"`
}