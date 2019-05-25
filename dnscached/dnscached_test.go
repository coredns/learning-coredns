package main

import (
	"bytes"
	"testing"
)

// TestCorefile tests the generation of the Corfile based on flags
func TestCorefile(t *testing.T) {
	resolv := []string{"/etc/resolv.conf"}
	tests := []struct {
		d        dnscached
		corefile string
	}{
		{dnscached{false, false, false, "127.0.0.1", 5300, 60, 10, 9984, 9984, resolv},
			".:5300 {\n errors\n bind 127.0.0.1\n cache 60 {\n  success 9984\n  denial 9984\n  prefetch 10\n }\n forward . /etc/resolv.conf \n}\n"},
		{dnscached{false, false, false, "127.0.0.1", 5300, 60, 10, 9984, 9984, []string{"1.1.1.1"}},
			".:5300 {\n errors\n bind 127.0.0.1\n cache 60 {\n  success 9984\n  denial 9984\n  prefetch 10\n }\n forward . 1.1.1.1 \n}\n"},
		{dnscached{false, false, false, "::1", 53, 120, 2, 50000, 9984, resolv},
			".:53 {\n errors\n bind ::1\n cache 120 {\n  success 50000\n  denial 9984\n  prefetch 2\n }\n forward . /etc/resolv.conf \n}\n"},
		{dnscached{false, false, false, "::1", 53, 120, 0, 50000, 9984, resolv},
			".:53 {\n errors\n bind ::1\n cache 120 {\n  success 50000\n  denial 9984\n }\n forward . /etc/resolv.conf \n}\n"},
		{dnscached{false, false, false, "::1", 53, 0, 2, 50000, 9984, resolv},
			".:53 {\n errors\n bind ::1\n forward . /etc/resolv.conf \n}\n"},
		{dnscached{false, false, true, "::1", 53, 0, 10, 50000, 9984, resolv},
			".:53 {\n errors\n bind ::1\n log\n forward . /etc/resolv.conf \n}\n"},
	}
	for i, test := range tests {
		input, err := test.d.corefile()

		if err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
			continue
		}

		text := bytes.NewBuffer(input.Body()).String()
		if test.corefile != text {
			t.Errorf("Test %d:\nexpected: %q\n   found: %q", i, test.corefile, text)
		}
	}
}
