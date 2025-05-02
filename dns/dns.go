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
	"fmt"
	"net"

	"github.com/soner3/net-scan/host"
)

type DnsResult struct {
	Host     string
	IPv4     string
	IPv6     string
	CNAME    string
	MX       string
	NS       string
	TXT      string
	NotFound bool
}

func NewDnsResult() *DnsResult {
	return &DnsResult{}
}

func lookupDns(host string) (*DnsResult, error) {
	res := NewDnsResult()

	cn, err := net.LookupCNAME(host)
	if err != nil {
		res.CNAME = "CNAME -\n"
	} else {
		res.CNAME = fmt.Sprintf("CNAME %s\n", cn)
	}

	mxs, err := net.LookupMX(host)
	if err != nil {
		res.MX = "MX -\n"
	} else {
		for _, mx := range mxs {
			res.MX = fmt.Sprintf("MX %s %d\n", mx.Host, mx.Pref)
		}
	}

	nss, err := net.LookupNS(host)
	if err != nil {
		res.NS = "Name Server: -\n"
	} else {
		for _, ns := range nss {
			res.NS = fmt.Sprintf("Name Server: %s\n", ns.Host)
		}
	}

	txts, err := net.LookupTXT(host)
	if err != nil {
		res.TXT = "TXT -\n"
	} else {
		for _, txt := range txts {
			res.TXT = fmt.Sprintf("TXT %s\n", txt)
		}
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		res.IPv4 = "IPv4 -\n"
		res.IPv6 = "IPv6 -\n"
	} else {
		for _, ip := range ips {
			if ip.To4() != nil {
				res.IPv4 = fmt.Sprintf("IPv4 %s\n", ip)
			} else {
				res.IPv6 = fmt.Sprintf("IPv6 %s\n", ip)
			}
		}

	}

	return res, nil
}

func Run(hl *host.HostList) ([]*DnsResult, error) {

	results := make([]*DnsResult, len(hl.Hosts))

	for i, h := range hl.Hosts {
		results[i] = NewDnsResult()
		results[i].Host = fmt.Sprintf("Host: %s\n", h)
		if _, err := net.LookupHost(h); err != nil {
			results[i].NotFound = true
			continue
		}
		var err error
		results[i], err = lookupDns(h)
		if err != nil {
			return nil, err
		}

	}
	return results, nil

}
