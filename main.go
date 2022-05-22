package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flagVersion := flag.Bool("version", false, "Show build information and exit")
	flagPullTime := flag.String("pull-time", "15s",
		"Time for network interface pulling")
	flagSystemdUnitName := flag.String("systemd-unit", "systemd-networkd.service",
		"Name of systemd unit to restart")
	flagNetworkInterfaceName := flag.String("interface-name", "wlan0",
		"Name of network interface to watch")

	flag.Parse()

	printInfo()

	if *flagVersion {
		return
	}

	pullTime, err := time.ParseDuration(*flagPullTime)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse 'pull-time' flag: %v\n", err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	for {
		select {
		case <-time.After(pullTime):
			if ok := linkHasRoutableAddressAndDefaultRoute(*flagNetworkInterfaceName); ok {
				continue
			}

			if err := tryRestartSystemdUnit(ctx, *flagSystemdUnitName); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)

				continue
			}

			_, _ = fmt.Fprintf(os.Stdout, "try-restart systemd unit '%s'\n", *flagSystemdUnitName)
		case <-ctx.Done():
			return
		}
	}
}
