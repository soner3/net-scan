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
package dns

import (
	"net"

	"github.com/soner3/net-scan/host"
)

type lookupResult struct {
	res DnsResult
}

type DnsResult struct {
	Host  string
	IPs   []net.IP
	CNAME string
	MXResult
	NSResult
	TXT      []string
	NotFound bool
}

type MXResult struct {
	NetMX []*net.MX
	MXErr error
}

type NSResult struct {
	NetNS []*net.NS
	NSErr error
}

func lookupDns(host string) DnsResult {
	res := DnsResult{}
	res.Host = host

	cn, err := net.LookupCNAME(host)
	if err != nil {
		res.NotFound = true
		return res
	} else {
		res.CNAME = cn
	}

	mxs, err := net.LookupMX(host)
	res.NetMX = mxs
	res.MXErr = err

	nss, err := net.LookupNS(host)
	res.NetNS = nss
	res.NSErr = err

	txts, err := net.LookupTXT(host)
	if err != nil {
		res.TXT = nil
	} else {
		res.TXT = txts
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		res.IPs = nil
	} else {
		res.IPs = ips
	}

	return res
}

func Run(hl *host.HostList) *[]DnsResult {
	results := make([]DnsResult, 0, len(hl.Hosts))
	res := make(chan DnsResult, len(hl.Hosts))

	for _, h := range hl.Hosts {
		h := h
		go func() {
			res <- lookupDns(h)
		}()
	}

	for range hl.Hosts {
		results = append(results, <-res)
	}

	return &results

}
