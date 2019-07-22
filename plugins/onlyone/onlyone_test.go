package onlyone

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func TestTrimRecords(t *testing.T) {
	baseAnswer := []dns.RR{
		test.CNAME("cname1.region2.skydns.test.	300	IN	CNAME		cname2.region2.skydns.test."),
		test.CNAME("cname6.region2.skydns.test.	300	IN	CNAME		endpoint.region2.skydns.test."),
		test.A("endpoint.region2.skydns.test.		300	IN	A			10.240.0.1"),
		test.A("endpoint.region2.skydns.test.		300	IN	A			10.240.0.1"),
		test.MX("mx.region2.skydns.test.			300	IN	MX		1	mx1.region2.skydns.test."),
		test.MX("mx.region2.skydns.test.			300	IN	MX		2	mx2.region2.skydns.test."),
		test.MX("mx.region2.skydns.test.			300	IN	MX		3	mx3.region2.skydns.test."),
		test.AAAA("endpoint.region2.skydns.test.	300	IN	AAAA		::1"),
		test.AAAA("endpoint.region2.skydns.test.	300	IN	AAAA		::2"),
	}

	tests := []struct {
		types  typeMap
		pick   func(int) int
		answer []dns.RR
	}{
		{
			types: typeMap{dns.TypeA: true, dns.TypeAAAA: true},
			pick:  func(int) int { return 0 },
			answer: []dns.RR{
				test.CNAME("cname1.region2.skydns.test.	300	IN	CNAME		cname2.region2.skydns.test."),
				test.CNAME("cname6.region2.skydns.test.	300	IN	CNAME		endpoint.region2.skydns.test."),
				test.A("endpoint.region2.skydns.test.		300	IN	A			10.240.0.1"),
				test.MX("mx.region2.skydns.test.			300	IN	MX		1	mx1.region2.skydns.test."),
				test.MX("mx.region2.skydns.test.			300	IN	MX		2	mx2.region2.skydns.test."),
				test.MX("mx.region2.skydns.test.			300	IN	MX		3	mx3.region2.skydns.test."),
				test.AAAA("endpoint.region2.skydns.test.	300	IN	AAAA		::1"),
			},
		},
	}

	for i, test := range tests {
		req := new(dns.Msg)
		req.SetQuestion("a.b.c.", dns.TypeA)
		req.Answer = baseAnswer

		o := &onlyone{types: test.types, pick: test.pick}
		o.trimRecords(req)

		if !sameAnswer(test.answer, req.Answer) {
			t.Errorf("Test %d: Expected %v, but got %v", i, test.answer, req.Answer)
		}
	}
}
