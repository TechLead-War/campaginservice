package utils

const (
	ErrMissingApp       = "missing app param"
	ErrMissingOS        = "missing os param"
	ErrMissingCountry   = "missing country param"
	ErrMethodNotAllowed = "method not allowed"
	InternalServerError = "internal server error"
)

var TargetingDimensions = []string{"app_id", "country", "os"}
