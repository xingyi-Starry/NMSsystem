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

func PendingError(t interface{}, err error) error {
	switch t := t.(type) {
	case Account:
		if t.Message == "Bad Gateway" || t.Message == "Internal Server Error" { // server crashed
			return NewNmsError("server crashed", ServerCrashErr)
		} else if t.Message != "" { // token expired
			return NewNmsError("token expired", TokenExpErr)
		}
	case TokenData:
		if t.Message == "Bad Gateway" || t.Message == "Internal Server Error" { // server crashed
			return NewNmsError("server crashed", ServerCrashErr)
		} else if t.Message != "" { // token expired
			return NewNmsError("token expired", TokenExpErr)
		}
	case Submission:
		if t.Message == "Bad Gateway" || t.Message == "Internal Server Error" { // server crashed
			return NewNmsError("server crashed", ServerCrashErr)
		} else if t.Message != "" { // submission expired
			return NewNmsError("submission expired", SubExpErr)
		}
	case SubResp:
		if t.Message == "Bad Gateway" || t.Message == "Internal Server Error" { // server crashed
			return NewNmsError("server crashed", ServerCrashErr)
		} else if t.Message != "" { // token expired
			return NewNmsError("token expired", TokenExpErr)
		}
	}
	return err
}
