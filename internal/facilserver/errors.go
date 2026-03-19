package facilserver

import "errors"

// Sentinel errors for facilitator API request validation.
var (
	ErrReadBody    = errors.New("failed to read request body")
	ErrInvalidJSON = errors.New("invalid JSON in request body")
)
