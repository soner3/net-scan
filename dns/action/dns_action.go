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
	"slices"
	"sort"

	"github.com/soner3/net-scan/dns"
	"github.com/soner3/net-scan/host"
)

type IndexOutput struct {
	Index  int
	Output string
}

func DnsAction(out io.Writer, filename string, search *[]string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}

	results := dns.Run(hl, search)
	labelWidth := 10

	outputChannel := make(chan IndexOutput, len(*results))
	defer close(outputChannel)

	for i, res := range *results {
		go func() {
			output := ""

			output += fmt.Sprintf("%s\n", res.Host)
			if res.NotFound {
				output += "\tNot Found\n\n"
				outputChannel <- IndexOutput{i, output}
				return
			}

			if slices.Contains(*search, "cname") {
				output += formatEntry("CNAME", res.CNAME, nil, labelWidth)
			}

			if slices.Contains(*search, "ip4") || slices.Contains(*search, "ip6") {
				for _, ip := range res.IPs {
					if ip.To4() != nil {
						if slices.Contains(*search, "ip4") {
							output += formatEntry("IPv4", ip.String(), nil, labelWidth)
						}
					} else {
						if slices.Contains(*search, "ip6") {
							output += formatEntry("IPv6", ip.String(), nil, labelWidth)
						}
					}
				}
			}

			if slices.Contains(*search, "ns") {
				for _, ns := range res.NetNS {
					output += formatEntry("NS", ns.Host, res.NSErr, labelWidth)
				}
			}

			if slices.Contains(*search, "mx") {
				for _, mx := range res.NetMX {
					output += formatEntry("MX", fmt.Sprintf("%s %d", mx.Host, mx.Pref), res.MXErr, labelWidth)
				}
			}

			if slices.Contains(*search, "txt") {
				for _, txt := range res.TXT {
					output += formatEntry("TXT", txt, nil, labelWidth)
				}
			}

			if output == fmt.Sprintf("%s\n", res.Host) {
				output += "\tNone\n"
			}

			output += "\n"
			outputChannel <- IndexOutput{i, output}
		}()
	}

	outputResult := ""
	collectedResults := make([]IndexOutput, len(*results))

	for i := range len(*results) {
		collectedResults[i] = <-outputChannel
	}

	sort.Slice(collectedResults, func(i, j int) bool {
		return collectedResults[i].Index < collectedResults[j].Index
	})

	for _, v := range collectedResults {
		outputResult += v.Output
	}

	fmt.Fprint(out, outputResult)
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
