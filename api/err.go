package api

const (
	TokenExpErr = iota
	SubExpErr
	ServerCrashErr
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
	return e.ErrType == TokenExpErr
}

func (e *NmsError) SubmissionExpired() bool {
	return e.ErrType == SubExpErr
}

func (e *NmsError) ServerCrashed() bool {
	return e.ErrType == ServerCrashErr
}
