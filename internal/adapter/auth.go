package adapter

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
)

type Authenticator interface {
	Authenticate(user, db, password string) error
}
type TrustAuth struct{}

func (TrustAuth) Authenticate(string, string, string) error { return nil }

type APIKeyAuth struct{ Key string }

var ErrInvalidCredentials = errors.New("invalid credentials")

func (a APIKeyAuth) Authenticate(user, db, password string) error {
	_ = user
	_ = db
	if a.Key == "" || subtle.ConstantTimeCompare([]byte(password), []byte(a.Key)) != 1 {
		return fmt.Errorf("invalid credentials: %v", ErrInvalidCredentials)
	}
	return nil
}
func BuildAuth(mode, key string) Authenticator {
	if strings.EqualFold(mode, "apikey") {
		return APIKeyAuth{Key: key}
	}
	return TrustAuth{}
}
