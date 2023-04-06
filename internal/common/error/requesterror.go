package error

// RequestError represent a request error
type RequestError struct {
	StatusCode int
	Err        error
}

// Error returns a string with the error
func (se *RequestError) Error() string {
	return se.Err.Error()
}

// StatusCode returns a StatusCode
func StatusCode(se RequestError) int {
	return se.StatusCode
}
