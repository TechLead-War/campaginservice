package utils

const (
	ErrMissingApp       = "missing app parameter"
	ErrMissingOS        = "missing os parameter"
	ErrMissingCountry   = "missing country parameter"
	ErrMethodNotAllowed = "method is not allowed"
	InternalServerError = "internal server error"
	DefaultApiPageLimit = 10
)

var TargetingDimensions = []string{"app_id", "country", "os"}
