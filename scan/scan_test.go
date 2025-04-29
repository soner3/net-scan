/*
Copyright © 2025 Soner Astan <sonerastan@icloud.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/scan"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name        string
		expectedOut string
	}{
		{"PortOpen", "open"},
		{"PortClosed", "closed"},
		{"PortTimeout", "timeout"},
	}

	localhost := "localhost"
	timeoutHost := "192.0.2.1"
	hl := host.NewHostList()
	hl.Add(localhost)
	ports := []int{}
	for _, tc := range testCases {
		if tc.name != "PortTimeout" {
			ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()
			_, portString, err := net.SplitHostPort(ln.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			port, err := strconv.Atoi(portString)
			if err != nil {
				t.Fatal(err)
			}
			ports = append(ports, port)

			if tc.name == "PortClosed" {
				ln.Close()
			}

		}

	}

	localhostRes := scan.Run(hl, &ports, "tcp", 1000)
	hl.Remove(localhost)
	hl.Add(timeoutHost)
	timeoutRes := scan.Run(hl, &ports, "tcp", 1000)

	if len(*(*localhostRes)[0].PortStates) != 2 {
		t.Errorf("Expected %d, got %d instead", 2, len(*localhostRes))
	}

	if len(*(*timeoutRes)[0].PortStates) != 2 {
		t.Errorf("Expected %d, got %d instead", 2, len(*timeoutRes))
	}

	if (*localhostRes)[0].Host != localhost {
		t.Errorf("Expected %s, got %s instead", localhost, (*localhostRes)[0].Host)
	}

	if (*timeoutRes)[0].Host != timeoutHost {
		t.Errorf("Expected %s, got %s instead", localhost, (*localhostRes)[0].Host)
	}

	if (*localhostRes)[0].NotFound {
		t.Errorf("Expected host %s to be found", (*localhostRes)[0].Host)
	}

	if (*timeoutRes)[0].NotFound {
		t.Errorf("Expected host %s to be found", (*localhostRes)[0].Host)
	}

	for _, res := range *localhostRes {
		if (*res.PortStates)[0].Open.String() != "open" {
			t.Errorf("Expected %s, got %s instead", "open", (*res.PortStates)[0].Open.String())
		}
		if (*res.PortStates)[1].Open.String() != "closed" {
			t.Errorf("Expected %s, got %s instead", "closed", (*res.PortStates)[1].Open.String())
		}
	}

	for _, res := range *timeoutRes {
		for _, ps := range *res.PortStates {
			if ps.Open.String() != "timeout" {
				t.Errorf("Expected %s, got %s instead", "timeout", ps.Open.String())
			}
		}
	}

}

func TestHostNotFound(t *testing.T) {
	testCases := []struct {
		hostname string
		found    bool
	}{
		{"localhost", true},
		{"unknown", false},
	}
	hl := host.NewHostList()
	ports := &[]int{80}

	for _, tc := range testCases {
		hl.Add(tc.hostname)
	}

	res := scan.Run(hl, ports, "tcp", 1000)

	for i, tc := range testCases {
		if (*res)[i].NotFound != !tc.found {
			t.Errorf("Expected %v, got %v instead", !tc.found, (*res)[i].NotFound)
		}
	}

}
