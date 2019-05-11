package onlyone

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type typeMap map[uint16]bool

type onlyone struct {
	Next  plugin.Handler
	zones []string
	types typeMap
}

func (o *onlyone) Name() string { return "onlyone" }

// ServeDNS implements the plugin.Handle interface.
func (o *onlyone) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// The request struct is a convenience struct
	state := request.Request{W: w, Req: r}

	// If the zone does not match one of ours, just pass it on
	if plugin.Zones(o.zones).Matches(state.Name()) == "" {
		return plugin.NextOrFailure(o.Name(), o.Next, ctx, w, r)
	}

	// The zone matches ours, so use a nonwriter to capture the response
	nw := nonwriter.New(w)
	status, err := plugin.NextOrFailure(o.Name(), o.Next, ctx, nw, r)

	// Now, examine the response and randomly choose a record to keep
	// TODO
	return status, err
}
