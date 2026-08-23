package adapter

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"time"
)

type TLSMode string

const (
	TLSDisabled TLSMode = "disabled"
	TLSOptional TLSMode = "optional"
	TLSRequired TLSMode = "required"
)

type TLSConfig struct {
	Mode                      TLSMode
	CertFile, KeyFile, CAFile string
	MinVersion                uint16
	HandshakeTimeout          time.Duration
}

func LoadTLS(c TLSConfig) (*tls.Config, error) {
	if c.Mode == TLSDisabled {
		return nil, nil
	}
	if c.CertFile == "" || c.KeyFile == "" {
		return nil, errors.New("certificate and key required")
	}
	cert, e := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if e != nil {
		return nil, e
	}
	out := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: c.MinVersion, ClientAuth: tls.NoClientCert}
	if c.CAFile != "" {
		b, e := os.ReadFile(c.CAFile)
		if e != nil {
			return nil, e
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(b) {
			return nil, errors.New("invalid CA")
		}
		out.ClientCAs = pool
		out.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return out, nil
}
