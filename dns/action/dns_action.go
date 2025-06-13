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
	"net"

	"github.com/soner3/net-scan/dns"
	"github.com/soner3/net-scan/host"
)

func DnsAction(out io.Writer, filename string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}

	result := dns.Run(hl)

	output := ""

	for _, res := range *result {
		output += fmt.Sprintf("%s\n", res.Host)
		if res.NotFound {
			output += "\tNot Found\n"
			continue
		}

		output += fmt.Sprintf("\tCNAME\t%s\n", res.CNAME)

		if res.IPs != nil {
			ipv4List := make([]net.IP, 0)
			ipv6List := make([]net.IP, 0)
			for _, ip := range *res.IPs {
				if ip.To4() != nil {
					ipv4List = append(ipv4List, ip)
				} else {
					ipv6List = append(ipv6List, ip)
				}

			}

			for _, ipv4 := range ipv4List {
				output += fmt.Sprintf("\tA\t%s\n", ipv4)
			}

			for _, ipv6 := range ipv6List {
				output += fmt.Sprintf("\tAAAA\t%s\n", ipv6)
			}
		}

		if res.MX != nil {
			for _, mx := range res.MX {
				output += fmt.Sprintf("\tMX\t%s\t%d\n", mx.Host, mx.Pref)
			}
		}

		if res.TXT != nil {
			for _, txt := range res.TXT {
				output += fmt.Sprintf("\tTXT\t%s\n", txt)
			}
		}

		if res.NS != nil {
			for _, ns := range res.NS {
				output += fmt.Sprintf("\tName Server\t%s\n", ns.Host)
			}
		}

		output += "\n"
	}

	fmt.Fprint(out, output)

	return nil
}
