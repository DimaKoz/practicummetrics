package error

type RequestError struct {
	StatusCode int
	Err        error
}

func (se *RequestError) Error() string {
	return se.Err.Error()
}

func StatusCode(se RequestError) int {
	return se.StatusCode
}
