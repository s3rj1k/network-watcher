package main

import (
	"context"
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
)

func tryRestartSystemdUnit(ctx context.Context, name string) error {
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %w", err)
	}

	defer conn.Close()

	ch := make(chan string)

	_, err = conn.TryRestartUnitContext(ctx, name, "replace", ch)
	if err != nil {
		return fmt.Errorf("failed to try-restart systemd unit '%s': %w", name, err)
	}

	<-ch

	return nil
}
