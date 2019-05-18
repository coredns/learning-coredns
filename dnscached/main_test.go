package main

import (
	"bytes"
	"testing"
)

// TestCorefile tests the generation of the Corfile based on flags
func TestCorefile(t *testing.T) {
	tests := []struct {
		bindIP                             string
		enableLog                          bool
		port, successSize, denialSize, ttl uint
		destinations                       []string
		shouldErr                          bool
		corefile                           string
	}{
		// oks
		{"", false, 5300, 9984, 9984, 60, nil, false,
			".:5300 {\n errors\n bind 127.0.0.1 ::1\n cache 60 {\n  success 9984\n  denial 9984\n  prefetch 10\n }\n forward . /etc/resolv.conf \n}\n"},
		{"::1", false, 53, 50000, 9984, 120, nil, false,
			".:53 {\n errors\n bind ::1\n cache 120 {\n  success 50000\n  denial 9984\n  prefetch 10\n }\n forward . /etc/resolv.conf \n}\n"},
		{"::1", false, 53, 50000, 9984, 0, nil, false,
			".:53 {\n errors\n bind ::1\n forward . /etc/resolv.conf \n}\n"},
		{"::1", true, 53, 50000, 9984, 0, nil, false,
			".:53 {\n errors\n bind ::1\n log\n forward . /etc/resolv.conf \n}\n"},
	}
	for i, test := range tests {
		bindIP = test.bindIP
		enableLog = test.enableLog
		port = test.port
		successSize = test.successSize
		denialSize = test.denialSize
		ttl = test.ttl
		input, err := corefile()

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

		text := bytes.NewBuffer(input.Body()).String()
		if test.corefile != text {
			t.Errorf("Test %d:\nexpected: %q\n   found: %q", i, test.corefile, text)
		}
	}
}
