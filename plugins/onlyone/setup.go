package onlyone

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/mholt/caddy"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("onlyone")

func init() {
	caddy.RegisterPlugin("onlyone", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	t, err := parse(c)
	if err != nil {
		return plugin.Error("onlyone", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		t.Next = next
		return t
	})

	return nil
}

func parse(c *caddy.Controller) (*onlyone, error) {
	o := &onlyone{types: typeMap{dns.TypeA: true, dns.TypeAAAA: true},
		pick: rand.Intn}

	found := false
	for c.Next() {
		// onlyone should just be in the server block once.
		if found {
			return nil, plugin.ErrOnce
		}
		found = true

		// parse the zone list, normalizing each to a FQDN, and
		// using the zones from the server block if none are given.
		args := c.RemainingArgs()
		if len(args) == 0 {
			o.zones = make([]string, len(c.ServerBlockKeys))
			copy(o.zones, c.ServerBlockKeys)
		}
		for _, str := range args {
			o.zones = append(o.zones, plugin.Host(str).Normalize())
		}

		for c.NextBlock() {
			switch c.Val() {
			case "types":
				args := c.RemainingArgs()
				if len(args) == 0 {
					return nil, errors.New(
						"at least one type must be listed")
				}
				o.types = make(typeMap, len(args))
				for _, a := range args {
					t, ok := dns.StringToType[strings.ToUpper(a)]
					if !ok {
						return nil,
							fmt.Errorf("invalid type %q",
								a)
					}
					o.types[t] = true
				}
			default:
				return nil, fmt.Errorf("invalid option %q", c.Val())
			}
		}
	}
	return o, nil
}
