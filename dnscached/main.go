package main

import (
	"fmt"
	"os"

	"github.com/coredns/caddy"
)

func init() {
	caddy.Quiet = true // don't show init stuff from caddy
	caddy.AppName = "dnscached"
	caddy.AppVersion = "1.0.1"
}

func main() {
	d := parseFlags()

	d.handleVersion()

	input, err := d.corefile()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	d.handleDryRun(input)

	// Start the server
	instance, err := caddy.Start(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Twiddle your thumbs
	instance.Wait()
}
