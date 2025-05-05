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
package action

import (
	"fmt"
	"io"

	"github.com/soner3/net-scan/dns"
	"github.com/soner3/net-scan/host"
)

func DnsAction(out io.Writer, filename string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}

	results := dns.Run(hl)
	output := ""

	labelWidth := 10

	for _, res := range *results {
		output += fmt.Sprintf("%s\n", res.Host)
		if res.NotFound {
			output += "\tNot Found\n\n"
			continue
		}

		output += formatEntry("CNAME", res.CNAME, nil, labelWidth)

		for _, ip := range res.IPs {
			if ip.To4() != nil {
				output += formatEntry("IPv4", ip.String(), nil, labelWidth)
			} else {
				output += formatEntry("IPv6", ip.String(), nil, labelWidth)
			}
		}

		for _, ns := range res.NetNS {
			output += formatEntry("NS", ns.Host, res.NSErr, labelWidth)
		}
		for _, mx := range res.NetMX {
			output += formatEntry("MX", fmt.Sprintf("%s %d", mx.Host, mx.Pref), res.MXErr, labelWidth)
		}
		for _, txt := range res.TXT {
			output += formatEntry("TXT", txt, nil, labelWidth)
		}

		output += "\n"
	}

	fmt.Fprint(out, output)
	return nil
}

func formatEntry(name string, value string, err error, width int) string {
	if value == "" {
		value = "-"
	}
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return fmt.Sprintf("\t%-*s %-50s %s\n", width, name, value, errStr)
}
