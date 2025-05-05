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
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/scan"
)

func TestRun(t *testing.T) {
	localhost := "localhost"
	timeoutHost := "192.0.2.1"
	hl := host.NewHostList()
	hl.Add(localhost)

	ports := []int{}
	portStatus := map[int]string{}

	ln1, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln1.Close()
	_, p1str, _ := net.SplitHostPort(ln1.Addr().String())
	p1, _ := strconv.Atoi(p1str)
	ports = append(ports, p1)
	portStatus[p1] = "open"

	ln2, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	_, p2str, _ := net.SplitHostPort(ln2.Addr().String())
	p2, _ := strconv.Atoi(p2str)
	ln2.Close()
	ports = append(ports, p2)
	portStatus[p2] = "closed"

	slices.Sort(ports)

	localhostRes := scan.Run(hl, &ports, "tcp", time.Second)

	hl.Remove(localhost)
	hl.Add(timeoutHost)
	timeoutRes := scan.Run(hl, &ports, "tcp", time.Second)

	findHost := func(results *[]scan.ScanResult, host string) *scan.ScanResult {
		for _, r := range *results {
			if r.Host == host {
				return &r
			}
		}
		return nil
	}

	res := findHost(localhostRes, localhost)
	if res == nil {
		t.Fatalf("Host %s not found in results", localhost)
	}
	if res.NotFound {
		t.Errorf("Expected host %s to be found", localhost)
	}
	if len(*res.PortStates) != 2 {
		t.Errorf("Expected 2 ports, got %d", len(*res.PortStates))
	}
	for _, ps := range *res.PortStates {
		want, ok := portStatus[ps.Port]
		if !ok {
			t.Errorf("Unexpected port: %d", ps.Port)
		}
		if ps.Open.String() != want {
			t.Errorf("Port %d: expected %s, got %s", ps.Port, want, ps.Open.String())
		}
	}

	res = findHost(timeoutRes, timeoutHost)
	if res == nil {
		t.Fatalf("Host %s not found in results", timeoutHost)
	}
	if res.NotFound {
		t.Errorf("Expected host %s to be found", timeoutHost)
	}
	for _, ps := range *res.PortStates {
		if ps.Open.String() != "timeout" {
			t.Errorf("Port %d: expected timeout, got %s", ps.Port, ps.Open.String())
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
