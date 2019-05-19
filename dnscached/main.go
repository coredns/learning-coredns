package main

import (
	"fmt"
	"os"

	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"

	"github.com/mholt/caddy"
)

const (
	AppVersion = "1.0.0"
	AppName    = "dnscached"
)

func init() {
	caddy.Quiet = true // don't show init stuff from caddy

	caddy.AppName = AppName
	caddy.AppVersion = AppVersion
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
