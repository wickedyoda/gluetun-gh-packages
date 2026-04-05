package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_apiKeyMethod_isAuthorized(t *testing.T) {
	t.Parallel()

	method := newAPIKeyMethod("secret-token")

	t.Run("header", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("X-API-Key", "secret-token")

		assert.True(t, method.isAuthorized(nil, request))
	})

	t.Run("query", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/?api_key=secret-token", nil)

		assert.True(t, method.isAuthorized(nil, request))
	})

	t.Run("mismatch", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("X-API-Key", "wrong-token")

		assert.False(t, method.isAuthorized(nil, request))
	})
}

func Test_basicAuthMethod_isAuthorized(t *testing.T) {
	t.Parallel()

	method := newBasicAuthMethod("alice", "super-secret")

	t.Run("authorized", func(t *testing.T) {
		t.Parallel()

		headers := make(http.Header)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetBasicAuth("alice", "super-secret")

		assert.True(t, method.isAuthorized(headers, request))
		assert.Empty(t, headers.Get("WWW-Authenticate"))
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		headers := make(http.Header)
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		assert.False(t, method.isAuthorized(headers, request))
		assert.Equal(t, `Basic realm="restricted", charset="UTF-8"`, headers.Get("WWW-Authenticate"))
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()

		headers := make(http.Header)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetBasicAuth("alice", "wrong")

		assert.False(t, method.isAuthorized(headers, request))
	})
}
