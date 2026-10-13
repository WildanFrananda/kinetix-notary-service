package notary

import (
	"errors"
	"fmt"
)

type Error struct {
	Kind    Kind
	Op      string
	Message string
	Err     error
}

func (e *Error) Error() string {
	text := fmt.Sprintf("%s: %s", e.Op, e.Kind)
	if e.Message != "" {
		text += ": " + e.Message
	}

	if e.Err != nil {
		text += fmt.Sprintf(": %v", e.Err)
	}

	return text
}

func (e *Error) Unwrap() error { return e.Err }

func Errf(kind Kind, op string, cause error, format string, args ...any) *Error {
	message := ""
	if format != "" {
		message = fmt.Sprintf(format, args...)
	}

	return &Error{Kind: kind, Op: op, Message: message, Err: cause}
}

func KindOf(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}

	return KindUnknown
}

func MessageOf(err error) string {
	for err != nil {
		var e *Error
		if !errors.As(err, &e) {
			return ""
		}

		if e.Message != "" {
			return e.Message
		}

		err = e.Err
	}

	return ""
}
