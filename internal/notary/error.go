package notary

import (
	"errors"
	"fmt"
)

type Error struct {
	Kind Kind
	Op   string
	Err  error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Op, e.Kind)
	}
	return fmt.Sprintf("%s: %s: %v", e.Op, e.Kind, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

func Errf(kind Kind, op string, cause error, format string, args ...any) *Error {
	if format == "" {
		return &Error{Kind: kind, Op: op, Err: cause}
	}
	if cause == nil {
		return &Error{Kind: kind, Op: op, Err: fmt.Errorf(format, args...)}
	}
	return &Error{Kind: kind, Op: op, Err: fmt.Errorf(format+": %w", append(args, cause)...)}
}

func KindOf(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return KindUnknown
}
