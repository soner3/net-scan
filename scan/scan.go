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
package scan

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/soner3/net-scan/host"
)

type PortState struct {
	Port    int
	Timeout bool
	Open    bool
}

type ScanResult struct {
	Host       string
	NotFound   bool
	PortStates *[]PortState
}

func NewScanResult(host string) *ScanResult {
	return &ScanResult{
		Host: host,
	}
}

func NewPortState(port int) *PortState {
	return &PortState{
		Port: port,
	}
}

func (ps PortState) String() string {
	if ps.Open {
		return "open"
	} else if ps.Timeout {
		return "timeout"
	} else {
		return "closed"
	}
}

// Scan the port on the given host
func scan(host string, port int, network string, timeout int) *PortState {
	ps := NewPortState(port)
	address := net.JoinHostPort(host, fmt.Sprintf("%d", ps.Port))
	con, err := net.DialTimeout(network, address, time.Millisecond*time.Duration(timeout))
	if err != nil {
		if os.IsTimeout(err) {
			ps.Timeout = true
		}
		return ps
	}
	con.Close()
	ps.Open = true
	return ps
}

// Run the scann process for all hosts
func Run(hl *host.HostList, ports *[]int, network string, timeout int) *[]ScanResult {
	results := make([]ScanResult, len(hl.Hosts))

	for i, h := range hl.Hosts {
		res := &results[i]
		res.Host = h
		res.PortStates = &[]PortState{}

		if _, err := net.LookupHost(h); err != nil {
			res.NotFound = true
			continue
		}

		for _, p := range *ports {
			ps := scan(h, p, network, timeout)
			*res.PortStates = append(*res.PortStates, *ps)
		}

	}

	return &results
}
