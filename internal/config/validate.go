package config

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func (c Config) Validate() error {
	if c.Listen == "" || c.ControlListen == "" {
		return errors.New("listen addresses required")
	}
	if _, e := net.ResolveTCPAddr("tcp", c.Listen); e != nil {
		return fmt.Errorf("proxy listen: %v", e)
	}
	if _, e := net.ResolveTCPAddr("tcp", c.ControlListen); e != nil {
		return fmt.Errorf("control listen: %w", e)
	}
	if c.MaxConnections < 1 {
		return errors.New("max connections must be positive")
	}
	if c.MaxMessageBytes < 64 {
		return errors.New("max message bytes too small")
	}
	if c.QueryTimeout <= 0 || c.IdleTransactionTimeout <= 0 {
		return errors.New("timeouts must be positive")
	}
	return nil
}
func (c Config) WithTimeout(d time.Duration) Config { c.QueryTimeout = d; return c }
