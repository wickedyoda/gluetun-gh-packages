package auth

import (
	"crypto/subtle"
	"net/http"
)

type basicAuthMethod struct {
	username []byte
	password []byte
}

func newBasicAuthMethod(username, password string) *basicAuthMethod {
	return &basicAuthMethod{
		username: []byte(username),
		password: []byte(password),
	}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (a *basicAuthMethod) equal(other authorizationChecker) bool {
	otherBasicMethod, ok := other.(*basicAuthMethod)
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare(a.username, otherBasicMethod.username) == 1 &&
		subtle.ConstantTimeCompare(a.password, otherBasicMethod.password) == 1
}

func (a *basicAuthMethod) isAuthorized(headers http.Header, request *http.Request) bool {
	username, password, ok := request.BasicAuth()
	if !ok {
		headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		return false
	}
	return subtle.ConstantTimeCompare([]byte(username), a.username) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), a.password) == 1
}
