package pgwire

import (
	"strconv"
	"strings"
)

type Parameters struct {
	User, Database, ApplicationName, ClientEncoding string
	Options                                         map[string]string
}

func ParseParameters(m map[string]string) Parameters {
	p := Parameters{Options: map[string]string{}}
	for k, v := range m {
		switch strings.ToLower(k) {
		case "user":
			p.User = v
		case "database":
			p.Database = v
		case "application_name":
			p.ApplicationName = v
		case "client_encoding":
			p.ClientEncoding = v
		default:
			p.Options[k] = v
		}
	}
	return p
}
func (p Parameters) String() string { return p.User + "/" + p.Database + "/" + p.ApplicationName }
func (p Parameters) Int(name string, def int) int {
	v := p.Options[name]
	n, e := strconv.Atoi(v)
	if e != nil {
		return def
	}
	return n
}
