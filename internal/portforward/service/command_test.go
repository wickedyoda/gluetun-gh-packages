//go:build linux

package service

import (
	"context"
	"testing"

	"github.com/qdm12/gluetun/internal/command"
	"github.com/stretchr/testify/require"
)

type testCmder struct {
	command string
}

type testLogger struct{}

func (testLogger) Debug(string) {}

func (testLogger) Info(string) {}

func (testLogger) Warn(string) {}

func (testLogger) Error(string) {}

func (c *testCmder) RunAndLog(_ context.Context, commandString string, _ command.Logger) (err error) {
	c.command = commandString
	return nil
}

func Test_Service_runCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cmder := &testCmder{}
	const commandTemplate = `/usr/bin/env echo {{PORTS}}-{{PORT}}-{{VPN_INTERFACE}}`
	ports := []uint16{1234, 5678}
	const vpnInterface = "tun0"
	logger := testLogger{}

	err := runCommand(ctx, cmder, logger, commandTemplate, ports, vpnInterface)

	require.NoError(t, err)
	require.Equal(t, "/usr/bin/env echo 1234,5678-1234-tun0", cmder.command)
}

func Test_validateCommandTemplate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args       []string
		errMessage string
	}{
		"valid command": {
			args: []string{"/usr/bin/env", "echo", "{{PORT}}"},
		},
		"empty": {
			errMessage: "command template is empty",
		},
		"templated executable blocked": {
			args:       []string{"/usr/bin/{{PORT}}"},
			errMessage: "template placeholders are not allowed in the executable path",
		},
		"shell command blocked": {
			args:       []string{"/bin/sh", "-c", "echo {{PORT}}"},
			errMessage: "shell command templates are not allowed: /bin/sh -c",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := validateCommandTemplate(testCase.args)

			if testCase.errMessage != "" {
				require.EqualError(t, err, testCase.errMessage)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
