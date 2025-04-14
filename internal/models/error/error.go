package error

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}
