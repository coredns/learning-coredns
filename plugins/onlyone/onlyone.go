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
	pick  func(int) int
}

func (o *onlyone) Name() string { return "onlyone" }

// ServeDNS implements the plugin.Handle interface.
func (o *onlyone) ServeDNS(ctx context.Context, w dns.ResponseWriter,
	r *dns.Msg) (int, error) {
	// The request struct is a convenience struct.
	state := request.Request{W: w, Req: r}

	// If the zone does not match one of ours, just pass it on.
	if plugin.Zones(o.zones).Matches(state.Name()) == "" {
		return plugin.NextOrFailure(o.Name(), o.Next, ctx, w, r)
	}

	// The zone matches ours, so use a nonwriter to capture the response.
	nw := nonwriter.New(w)

	// Call all the next plugin in the chain.
	rcode, err := plugin.NextOrFailure(o.Name(), o.Next, ctx, nw, r)
	if err != nil {
		// Simply return if there was an error.
		return rcode, err
	}

	// Now we know that a successful response was received from a plugin
	// that appears later in the chain. Next is to examine that response
	// and trim out extra records, then write it to the client.
	w.WriteMsg(o.trimRecords(nw.Msg))
	return rcode, err
}

func (o *onlyone) trimRecords(m *dns.Msg) *dns.Msg {
	// The trimming behavior is relatively expensive, so if there is one
	// or fewer answers, we know it doesn't apply so just return.
	if len(m.Answer) <= 1 {
		return m
	}

	// Allocate an array to hold answers to keep.
	keep := make([]bool, len(m.Answer))

	// Allocate a map to correlate each subject type to a list of indexes.
	indexes := make(map[uint16][]int, len(o.types)/2)

	// Loop through the answers, either deciding to keep it, or putting
	// it in a provisional list of indexes for a subject type.
	for i, a := range m.Answer {
		h := a.Header()
		if _, ok := o.types[h.Rrtype]; ok {
			// this type is subject to this plugin, so stash
			// away the index of this record for later.
			provisional, _ := indexes[h.Rrtype]
			indexes[h.Rrtype] = append(provisional, i)
		} else {
			// not subject to this plugin, so we keep it.
			keep[i] = true
		}
	}

	// Now we loop through each type with multiple records and pick one.
	for _, provisional := range indexes {
		keep[provisional[o.pick(len(provisional))]] = true
	}

	// Now copy the ones we want to keep into a new Answer list.
	var newAnswer []dns.RR
	for i, a := range m.Answer {
		if keep[i] {
			newAnswer = append(newAnswer, a)
		}
	}
	m.Answer = newAnswer
	return m
}
