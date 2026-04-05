package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/qdm12/gluetun/internal/command"
)

var (
	ErrCommandTemplateEmpty        = errors.New("command template is empty")
	ErrShellCommandTemplateBlocked = errors.New("shell command templates are not allowed")
	ErrExecutableTemplateBlocked   = errors.New("template placeholders are not allowed in the executable path")
)

func runCommand(ctx context.Context, cmder Cmder, logger Logger,
	commandTemplate string, ports []uint16, vpnInterface string,
) (err error) {
	templateArgs, err := command.Split(commandTemplate)
	if err != nil {
		return fmt.Errorf("parsing command template: %w", err)
	}
	err = validateCommandTemplate(templateArgs)
	if err != nil {
		return err
	}

	portStrings := make([]string, len(ports))
	for i, port := range ports {
		portStrings[i] = fmt.Sprint(int(port))
	}
	portsString := strings.Join(portStrings, ",")

	commandString := strings.ReplaceAll(commandTemplate, "{{PORTS}}", portsString)
	commandString = strings.ReplaceAll(commandString, "{{PORT}}", portStrings[0])
	commandString = strings.ReplaceAll(commandString, "{{VPN_INTERFACE}}", vpnInterface)
	return cmder.RunAndLog(ctx, commandString, logger)
}

func validateCommandTemplate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%w", ErrCommandTemplateEmpty)
	}

	if strings.Contains(args[0], "{{") {
		return fmt.Errorf("%w", ErrExecutableTemplateBlocked)
	}

	shell := filepath.Base(args[0])
	switch shell {
	case "sh", "ash", "bash", "dash", "ksh", "zsh":
		if len(args) >= 2 && strings.HasPrefix(args[1], "-c") {
			return fmt.Errorf("%w: %s %s", ErrShellCommandTemplateBlocked, args[0], args[1])
		}
	}

	return nil
}
