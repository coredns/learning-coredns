package onlyone

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"

	"github.com/miekg/dns"
)

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
	o := &onlyone{types: typeMap{dns.TypeA: true, dns.TypeAAAA: true}}

	found := false
	for c.Next() {
		// onlyone should just be in the server block once
		if found {
			return nil, plugin.ErrOnce
		}
		found = true
		args := c.RemainingArgs()
		if len(args) == 0 {
			o.zones = []string{"."} // match any zone
		} else {
			o.zones = args
		}
		for c.NextBlock() {
			switch c.Val() {
			case "types":
				args := c.RemainingArgs()
				if len(args) == 0 {
					return nil, errors.New("at least one type must be listed with types")
				}
				o.types = make(typeMap, len(args))
				for _, a := range args {
					t, ok := dns.StringToType[strings.ToUpper(a)]
					if !ok {
						return nil, fmt.Errorf("%s is not a valid type", a)
					}
					o.types[t] = true
				}
			default:
				return nil, fmt.Errorf("%s is an invalid option", c.Val())
			}
		}
	}
	return o, nil
}
