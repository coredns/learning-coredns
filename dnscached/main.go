package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	//"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/mholt/caddy"
)

var (
	version, enableLog, dryRun         bool
	bindIP                             string
	port, ttl, successSize, denialSize uint
	destinations                       []string
)

const (
	prefetchAmount = 10
	AppVersion     = "1.0.0"
	AppName        = "dnscached"
	defaultDest    = "/etc/resolv.conf"
)

func init() {
	caddy.Quiet = true // don't show init stuff from caddy

	caddy.AppName = AppName
	caddy.AppVersion = AppVersion
}

func setupFlags() *flag.FlagSet {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&caddy.PidFile, "pidfile", "", "File `path` to write pid file")
	f.BoolVar(&version, "version", false, "Show version")
	f.BoolVar(&dryRun, "dry-run", false, "Prints out the internally generated Corefile and exits")
	f.BoolVar(&enableLog, "log", false, "Enable query logging")
	f.StringVar(&bindIP, "bind", "", "`IP(s)` to which to bind (default '127.0.0.1 ::1')")
	f.UintVar(&port, "port", 5300, "Local port `number` to use")
	f.UintVar(&successSize, "success", 9984, "Number of success cache `entries`")
	f.UintVar(&denialSize, "denial", 9984, "Number of denial cache `entries`")
	f.UintVar(&ttl, "ttl", 60, "Maximum `seconds` to cache records, zero disables caching")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE\n-----\n%s [ options ] [ destinations ]\n", os.Args[0])
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

	return f
}

func main() {
	caddy.TrapSignals()

	flag.CommandLine = setupFlags()
	flag.Parse()

	if bindIP == "" {
		bindIP = "127.0.0.1 ::1"
	}
	destinations = flag.Args()
	if len(destinations) == 0 {
		destinations = []string{"/etc/resolv.conf"}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Set to 0 because we're doing our own time, with timezone

	if version {
		showVersion()
		os.Exit(0)
	}

	input, err := corefile()
	if err != nil {
		log.Fatal(err)
	}

	if dryRun {
		fmt.Print(bytes.NewBuffer(input.Body()).String())
		os.Exit(0)
	}

	showVersion()

	// Start the server
	instance, err := caddy.Start(input)
	if err != nil {
		log.Fatal(err)
	}


	// Twiddle your thumbs
	instance.Wait()
}

// corefile generates the Corefile based on the flags
func corefile() (caddy.Input, error) {

	var b bytes.Buffer
	_, err := b.WriteString(fmt.Sprintf(".:%d {\n errors\n bind %s\n", port, bindIP))
	if err != nil {
		return nil, err
	}

	if enableLog {
		_, err = b.WriteString( " log\n")
		if err != nil {
			return nil, err
		}
	}
	if ttl > 0 {
		_, err = b.WriteString(fmt.Sprintf(" cache %d {\n  success %d\n  denial %d\n  prefetch %d\n }\n",
			ttl, successSize, denialSize, prefetchAmount))
		if err != nil {
			return nil, err
		}
	}

	_, err = b.WriteString(" forward . ")
	if err != nil {
		return nil, err
	}
	for _, d := range destinations {
		_, err = b.WriteString(d)
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

// logVersion logs the version that is starting.
func logVersion() {
	clog.Info(versionString())
}

// showVersion prints the version that is starting.
func showVersion() {
	fmt.Print(versionString())
}

// versionString returns the CoreDNS version as a string.
func versionString() string {
	return fmt.Sprintf("%s-%s\n", caddy.AppName, caddy.AppVersion)
}

// flagsBlacklist removes flags with these names from our flagset.
var flagsBlacklist = map[string]struct{}{
	"logtostderr":      {},
	"alsologtostderr":  {},
	"v":                {},
	"stderrthreshold":  {},
	"vmodule":          {},
	"log_backtrace_at": {},
	"log_dir":          {},
}

var flagsToKeep []*flag.Flag
