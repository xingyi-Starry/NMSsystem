package api

const (
	tokenExpErr = iota
	subExpErr
	serverCrashErr
)

type NmsError struct {
	err     string
	ErrType int
}

func (e *NmsError) Error() string {
	return e.err
}

func NewNmsError(err string, errType int) *NmsError {
	return &NmsError{err: err, ErrType: errType}
}

func (e *NmsError) TokenExpired() bool {
	return e.ErrType == tokenExpErr
}

func (e *NmsError) SubmissionExpired() bool {
	return e.ErrType == subExpErr
}

func (e *NmsError) ServerCrashed() bool {
	return e.ErrType == serverCrashErr
}
