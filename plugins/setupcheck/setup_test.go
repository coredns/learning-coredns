package setupcheck

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`setupcheck a {
				foo
				bar
				foobar foo bar
		}
		setupcheck b {
			bar bar bar
		}
		setupcheck c
		setupcheck d {
		}`},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		err := parse(c)
		if err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
			continue
		}
	}
}
