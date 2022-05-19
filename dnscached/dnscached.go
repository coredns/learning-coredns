package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"

	"github.com/coredns/caddy"
)

const (
	defaultDest = "/etc/resolv.conf"
)

type dnscached struct {
	printVersion, dryRun, enableLog bool
	bindIP                          string
	port, ttl, prefetchAmount       uint
	successSize, denialSize         uint
	destinations                    []string
}

func parseFlags() *dnscached {
	d := &dnscached{}
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&caddy.PidFile, "pidfile", "", "File `path` to write pid file")
	f.BoolVar(&d.printVersion, "version", false, "Show version")
	f.BoolVar(&d.dryRun, "dry-run", false,
		"Prints out the internally generated Corefile and exits")
	f.BoolVar(&d.enableLog, "log", false, "Enable query logging")
	f.StringVar(&d.bindIP, "bind", "127.0.0.1 ::1", "`IP(s)` to which to bind")
	f.UintVar(&d.port, "port", 5300, "Local port `number` to use")
	f.UintVar(&d.successSize, "success", 9984,
		"Number of success cache `entries`")
	f.UintVar(&d.denialSize, "denial", 9984, "Number of denial cache `entries`")
	f.UintVar(&d.prefetchAmount, "prefetch",
		10, "Times a query must be made per minute to qualify for prefetch")
	f.UintVar(&d.ttl, "ttl", 60,
		"Maximum `seconds` to cache records, zero disables caching")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"USAGE\n-----\n%s [ options ] [ destinations ]\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOPTIONS\n-------\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDESTINATIONS\n------------")
		fmt.Fprintf(os.Stderr, `
One or more forwarding destinations. Each can be a file in /etc/resolv.conf
format or a destination IP or IP:PORT, with or without a with or without a
protocol (leading "PROTO://"), with "dns" and "tls" as the supported PROTO
values. If omitted, "dns" is assumed as the protocol. The default destination is
/etc/resolv.conf.
`)
	}

	flag.CommandLine = f
	flag.Parse()
	d.destinations = flag.Args()
	if len(d.destinations) == 0 {
		d.destinations = []string{"/etc/resolv.conf"}
	}

	return d
}

func (d *dnscached) handleVersion() {
	if d.printVersion {
		fmt.Printf("%s-%s\n", caddy.AppName, caddy.AppVersion)
		os.Exit(0)
	}
}

func (d *dnscached) handleDryRun(input caddy.Input) {
	if d.dryRun {
		fmt.Print(bytes.NewBuffer(input.Body()).String())
		os.Exit(0)
	}
}

// corefile generates the Corefile based on the flags
func (d *dnscached) corefile() (caddy.Input, error) {
	var b bytes.Buffer
	_, err := b.WriteString(fmt.Sprintf(".:%d {\n errors\n bind %s\n",
		d.port, d.bindIP))
	if err != nil {
		return nil, err
	}

	if d.enableLog {
		_, err = b.WriteString(" log\n")
		if err != nil {
			return nil, err
		}
	}
	if d.ttl > 0 {
		_, err = b.WriteString(fmt.Sprintf(" cache %d {\n  success %d\n  denial %d\n",
			d.ttl, d.successSize, d.denialSize))
		if err != nil {
			return nil, err
		}
		if d.prefetchAmount > 0 {
			_, err = b.WriteString(fmt.Sprintf("  prefetch %d\n", d.prefetchAmount))
			if err != nil {
				return nil, err
			}
		}
		_, err = b.WriteString(" }\n")
		if err != nil {
			return nil, err
		}
	}

	_, err = b.WriteString(" forward . ")
	if err != nil {
		return nil, err
	}
	for _, dest := range d.destinations {
		_, err = b.WriteString(dest)
		if err != nil {
			return nil, err
		}
		_, err = b.WriteString(" ")
		if err != nil {
			return nil, err
		}
	}
	_, err = b.WriteString("\n}\n")
	if err != nil {
		return nil, err
	}

	return caddy.CaddyfileInput{
		Contents:       b.Bytes(),
		Filepath:       "<flags>",
		ServerTypeName: "dns",
	}, nil
}
