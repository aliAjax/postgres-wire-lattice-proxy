package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Listen, ControlListen, BackendAddr, AuthMode string
	MaxConnections, MaxMessageBytes              int
	QueryTimeout, IdleTransactionTimeout         time.Duration
	APIKey                                       string
}

func Default() Config {
	return Config{Listen: "127.0.0.1:6432", ControlListen: "127.0.0.1:8080", BackendAddr: "simulator", AuthMode: "trust", MaxConnections: 128, MaxMessageBytes: 1 << 20, QueryTimeout: 10 * time.Second, IdleTransactionTimeout: time.Minute}
}
func Load(path string) (Config, error) {
	c := Default()
	if path != "" {
		if _, err := os.Stat(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return c, err
		}
	}
	env(&c)
	return c, nil
}
func env(c *Config) {
	if v := os.Getenv("PROXY_LISTEN"); v != "" {
		c.Listen = v
	}
	if v := os.Getenv("CONTROL_LISTEN"); v != "" {
		c.ControlListen = v
	}
	if v := os.Getenv("BACKEND_ADDR"); v != "" {
		c.BackendAddr = v
	}
	if v := os.Getenv("AUTH_MODE"); v != "" {
		c.AuthMode = v
	}
	if v := os.Getenv("API_KEY"); v != "" {
		c.APIKey = v
	}
	if v := os.Getenv("MAX_CONNECTIONS"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.MaxConnections = n
		} else {
			c.MaxConnections = 0
		}
	}
}
