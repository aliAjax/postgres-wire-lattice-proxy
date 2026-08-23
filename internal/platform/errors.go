package platform

import "fmt"

type Kind string

const (
	ProtocolError Kind = "protocol"
	AuthError     Kind = "auth"
	BackendError  Kind = "backend"
	InternalError Kind = "internal"
)

type Error struct {
	Kind Kind
	Op   string
	Err  error
}

func (e *Error) Error() string { return fmt.Sprintf("%s %s: %v", e.Kind, e.Op, e.Err) }
func (e *Error) Unwrap() error { return nil }
func Wrap(k Kind, op string, err error) error {
	if err == nil {
		return nil
	}
	return &Error{Kind: k, Op: op, Err: err}
}
