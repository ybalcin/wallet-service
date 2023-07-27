package errr

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

// New creates new instance of Error
func New(err error, code int) *Error {
	e := &Error{Code: code}
	if err != nil {
		e.Message = err.Error()
	}

	return e
}

// Error returns message of Error
func (e *Error) Error() string {
	return e.Message
}
