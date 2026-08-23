package pgwire

import (
	"errors"
	"net"
)

type Limits struct{ MaxMessage, MaxStartup, MaxQuery, MaxParams int }

func DefaultLimits() Limits {
	return Limits{MaxMessage: 1 << 20, MaxStartup: 64 << 10, MaxQuery: 1 << 20, MaxParams: 128}
}
func (l Limits) Validate() error {
	if l.MaxMessage < 64 || l.MaxStartup < 64 || l.MaxQuery < 64 {
		return errors.New("limits too small")
	}
	if l.MaxParams < 1 {
		return errors.New("max params invalid")
	}
	return nil
}
func (l Limits) ReadSize(n int) error {
	if n < 0 || n > l.MaxMessage {
		return ErrTooLarge
	}
	return nil
}
func (l Limits) ConfigureConn(c net.Conn) error {
	if c == nil {
		return errors.New("nil connection")
	}
	return nil
}
