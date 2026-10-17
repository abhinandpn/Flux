package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/abhinandpn/flux/internal/core"
	"github.com/abhinandpn/flux/internal/monitor"
	"github.com/abhinandpn/flux/internal/platform/darwin"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "flux:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "watch":
		return watch(args[1:])

	case "status":
		return status()

	case "interfaces":
		return interfaces()

	case "help", "--help", "-h":
		printUsage()
		return nil

	default:
		printUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func watch(args []string) error {
	flags := flag.NewFlagSet("watch", flag.ContinueOnError)

	interfaceName := flags.String(
		"interface",
		"en0",
		"Network interface to monitor, e.g. en0",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() > 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	fmt.Println("Flux — Your network. At a glance.")
	fmt.Println("Monitoring network traffic. Press Ctrl+C to stop.")

	updates := make(chan monitor.Update)
	done := make(chan error, 1)

	m := monitor.Monitor{
		Sampler: monitor.Sampler{
			Provider: darwin.Provider{},
			Interval: time.Second,
		},
		Interface: *interfaceName,
	}

	go func() {
		done <- m.Run(ctx, updates)
	}()

	for {
		select {
		case <-ctx.Done():
			<-done
			fmt.Println("\nFlux stopped.")
			return nil

		case err := <-done:
			return err

		case update := <-updates:
			if update.Err != nil {
				fmt.Fprintf(
					os.Stderr,
					"\nMonitoring: %v\n",
					update.Err,
				)
				continue
			}

			fmt.Printf(
				"\r%-10s ↓ %-14s ↑ %-14s",
				update.Rate.Interface,
				core.FormatRate(update.Rate.DownloadBPS),
				core.FormatRate(update.Rate.UploadBPS),
			)
		}
	}
}

func status() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	samples, err := (darwin.Provider{}).ReadCounters(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Flux status: ready")
	fmt.Printf("Interfaces with counters: %d\n", len(samples))
	fmt.Println("Live sampling starts with: flux watch")

	return nil
}

func interfaces() error {
	items, err := darwin.ListInterfaces()
	if err != nil {
		return err
	}

	fmt.Println("Network interfaces:")

	for _, item := range items {
		state := "down"
		if item.Up {
			state = "up"
		}

		if item.Loopback {
			state += ", loopback"
		}

		fmt.Printf(
			"  %-12s index=%-3d %s\n",
			item.Name,
			item.Index,
			state,
		)
	}

	return nil
}

func printUsage() {
	fmt.Println(`Flux — Your network. At a glance.

Usage:
  flux watch [--interface NAME]  Monitor live network traffic
  flux status                    Check counter availability
  flux interfaces                List network interfaces
  flux help                      Show this help`)
}
