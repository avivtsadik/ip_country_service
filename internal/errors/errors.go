package errors

import "errors"

// DataStore errors
var (
	ErrIPNotFound      = errors.New("IP address not found")
	ErrInvalidIP       = errors.New("invalid IP address format")
	ErrDatastoreInit   = errors.New("failed to initialize datastore")
)

// Rate limiter errors
var (
	ErrRateLimited = errors.New("rate limit exceeded")
)

// HTTP errors
var (
	ErrMissingIPParam = errors.New("missing ip parameter")
	ErrInternalServer = errors.New("internal server error")
)