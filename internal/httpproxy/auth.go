package httpproxy

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func (h *handler) isAuthorized(responseWriter http.ResponseWriter, request *http.Request) (authorized bool) {
	if h.username == "" || (request.Method != http.MethodConnect && !request.URL.IsAbs()) {
		return true
	}
	basicAuth := request.Header.Get("Proxy-Authorization")
	if basicAuth == "" {
		responseWriter.Header().Set("Proxy-Authenticate", `Basic realm="Access to Gluetun over HTTP"`)
		responseWriter.WriteHeader(http.StatusProxyAuthRequired)
		return false
	}
	b64UsernamePassword := strings.TrimPrefix(basicAuth, "Basic ")
	b, err := base64.StdEncoding.DecodeString(b64UsernamePassword)
	if err != nil {
		h.logger.Info("Cannot decode Proxy-Authorization header value from " +
			request.RemoteAddr + ": " + err.Error())
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return false
	}
	usernamePassword := strings.Split(string(b), ":")
	const expectedFields = 2
	if len(usernamePassword) != expectedFields {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return false
	}
	if h.username != usernamePassword[0] || h.password != usernamePassword[1] {
		h.logger.Info("proxy authentication mismatch from " + request.RemoteAddr)
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}
