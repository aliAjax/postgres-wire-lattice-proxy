package frontend

import (
	"crypto/tls"
	"net"
)

type TLSListener struct {
	net.Listener
	Config *tls.Config
}

func (l TLSListener) AcceptTLS() (net.Conn, error) {
	c, e := l.Listener.Accept()
	if e != nil {
		return nil, e
	}
	t := tls.Server(c, l.Config)
	if e = t.Handshake(); e != nil {
		_ = c.Close()
		return nil, e
	}
	return t, nil
}
