package onlyone

import (
	"testing"

	"github.com/mholt/caddy"

	"github.com/miekg/dns"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		types     typeMap
		zones     []string
	}{
		// oks
		{`onlyone`, false, typeMap{dns.TypeA: true, dns.TypeAAAA: true}, []string{"."}},
		{`onlyone foo.com`, false, typeMap{dns.TypeA: true, dns.TypeAAAA: true}, []string{"foo.com"}},
		{`onlyone foo.com example.org`, false,
			typeMap{dns.TypeA: true, dns.TypeAAAA: true}, []string{"foo.com", "example.org"}},
		{"onlyone {\n types a\n}", false, typeMap{dns.TypeA: true}, []string{"."}},
		{"onlyone {\n types a mx aaaa hinfo\n}", false,
			typeMap{dns.TypeA: true, dns.TypeMX: true, dns.TypeAAAA: true, dns.TypeHINFO: true},
			[]string{"."}},
		// fails
		{"onlyone {\n types foo\n}", true, nil, []string{"."}},
		{"onlyone\nonlyone\n", true, nil, []string{"."}},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		p, err := parse(c)
		if test.shouldErr && err == nil {
			t.Errorf("Test %v: Expected error but found nil", i)
			continue
		} else if !test.shouldErr && err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
			continue
		}

		if test.shouldErr {
			continue
		}

		if !typeMapsMatch(test.types, p.types) {
			t.Errorf("Test %v: Expected types %v but found: %v", i, test.types, p.types)
		}
	}
}

func typeMapsMatch(a, b typeMap) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || v != bv {
			return false
		}
	}
	return true
}
