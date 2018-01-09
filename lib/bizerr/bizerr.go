/*Package bizerr provides business errors that can signal us which response code
should be sent to the user.
*/
package bizerr

import (
	"github.com/pkg/errors"
)

// ErrorType is an enum type for error codes.
type ErrorType uint

const (
	// ErrorUnknown is an error that has not been annotated.
	ErrorUnknown ErrorType = iota
	// ErrorUserInput is an error that happened because of incorrect user input.
	ErrorUserInput
	// ErrorUnauthorized is returned when user isn't allowed to do some action.
	ErrorUnauthorized
	// ErrorNotFound is returned when a requested object couldn't be found.
	ErrorNotFound
	// ErrorConflict is returned when there's a conflict between request(s) and our data.
	ErrorConflict
)

type bizErr struct {
	err     error
	errType ErrorType
}

func (e bizErr) Cause() error {
	return e.err
}

func (e bizErr) Error() string {
	return e.err.Error()
}

func (e bizErr) ErrorType() ErrorType {
	return e.errType
}

// Wrap annotates an error with ErrorType.
func Wrap(err error, t ErrorType) error {
	return bizErr{err: err, errType: t}
}

// New creates a new error with preset Type.
func New(err string, t ErrorType) error {
	return bizErr{err: errors.New(err), errType: t}
}

// Type returns an ErrorType of an error.
func Type(err error) ErrorType {
	type causer interface {
		Cause() error
	}

	type ErrorTyper interface {
		ErrorType() ErrorType
	}

	for err != nil {
		typed, ok := err.(ErrorTyper)
		if ok {
			return typed.ErrorType()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return ErrorUnknown
}
