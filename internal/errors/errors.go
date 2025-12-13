package errors

import "errors"

// DataStore errors
var (
	ErrIPNotFound               = errors.New("IP address not found")
	ErrInvalidIP                = errors.New("invalid IP address format")
	ErrUnsupportedDatastoreType = errors.New("unsupported datastore type")
	ErrDatastoreLookupFailed    = errors.New("datastore lookup failed")
)

// ErrRateLimited Rate limiter errors
var (
	ErrRateLimited = errors.New("rate limit exceeded")
)

// HTTP errors
var (
	ErrMissingIPParam   = errors.New("missing ip parameter")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInternalServer   = errors.New("internal server error")
)

// ErrAppInit Application errors
var (
	ErrAppInit = errors.New("failed to initialize application")
)
