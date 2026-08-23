package adapter

import (
	"crypto/subtle"
	"errors"
	"strings"
)

type Authenticator interface {
	Authenticate(user, db, password string) error
}
type TrustAuth struct{}

func (TrustAuth) Authenticate(string, string, string) error { return nil }

type APIKeyAuth struct{ Key string }

func (a APIKeyAuth) Authenticate(user, db, password string) error {
	_ = user
	_ = db
	if a.Key == "" || subtle.ConstantTimeCompare([]byte(password), []byte(a.Key)) != 1 {
		return errors.New("invalid credentials")
	}
	return nil
}
func BuildAuth(mode, key string) Authenticator {
	if strings.EqualFold(mode, "apikey") {
		return APIKeyAuth{Key: key}
	}
	return TrustAuth{}
}
