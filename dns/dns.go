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

type DnsResult struct {
	Host     string
	IPs      *[]net.IP
	CNAME    string
	MX       []*net.MX
	NS       []*net.NS
	TXT      []string
	NotFound bool
}

func lookupDns(host string) DnsResult {
	res := DnsResult{Host: host}

	cn, err := net.LookupCNAME(host)
	if err != nil {
		res.NotFound = true
		return res
	} else {
		res.CNAME = cn
	}

	mxs, err := net.LookupMX(host)
	if err != nil {
		res.MX = nil
	} else {
		res.MX = mxs
	}

	nss, err := net.LookupNS(host)
	if err != nil {
		res.NS = nil
	} else {
		res.NS = nss
	}

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
		res.IPs = &ips
	}

	return res
}

func Run(hl *host.HostList) *[]DnsResult {
	results := make([]DnsResult, len(hl.Hosts))

	for i, h := range hl.Hosts {
		results[i] = lookupDns(h)
	}

	return &results

}
