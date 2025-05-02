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

	result, err := dns.Run(hl)
	if err != nil {
		return err
	}

	output := ""

	for _, res := range result {
		output += res.Host
		if res.NotFound {
			output += fmt.Sprint(out, "Not Found\n")
			continue
		}

		if res.CNAME != "" {
			output += res.CNAME
		}
		if res.NS != "" {
			output += res.NS
		}
		if res.IPv4 != "" {
			output += res.IPv4
		} else {
			output += "-\n"
		}
		if res.IPv6 != "" {
			output += res.IPv6
		} else {
			output += "-\n"
		}
		if res.MX != "" {
			output += res.MX
		}
		if res.TXT != "" {
			output += res.TXT
		}
		output += "\n"
	}

	fmt.Fprint(out, output)

	return nil
}
