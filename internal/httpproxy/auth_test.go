package httpproxy

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	infoMessages  []string
	errorMessages []string
	debugMessages []string
}

func (l *testLogger) Info(s string)  { l.infoMessages = append(l.infoMessages, s) }
func (l *testLogger) Error(s string) { l.errorMessages = append(l.errorMessages, s) }
func (l *testLogger) Debug(s string) { l.debugMessages = append(l.debugMessages, s) }
func (l *testLogger) Warn(string)    {}

func Test_handler_isAuthorized_doesNotLogCredentials(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	handler := &handler{
		username: "alice",
		password: "very-secret-password",
		logger:   logger,
	}

	request := httptest.NewRequest(http.MethodConnect, "http://example.com", nil)
	request.RemoteAddr = "10.0.0.2:12345"
	encoded := base64.StdEncoding.EncodeToString([]byte("alice:wrong-password"))
	request.Header.Set("Proxy-Authorization", "Basic "+encoded)

	recorder := httptest.NewRecorder()
	authorized := handler.isAuthorized(recorder, request)

	assert.False(t, authorized)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, []string{"proxy authentication mismatch from 10.0.0.2:12345"}, logger.infoMessages)

	logOutput := strings.Join(append(append([]string{}, logger.infoMessages...), logger.debugMessages...), "\n")
	assert.NotContains(t, logOutput, "alice")
	assert.NotContains(t, logOutput, "wrong-password")
	assert.NotContains(t, logOutput, "very-secret-password")
}
