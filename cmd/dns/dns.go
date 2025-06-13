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
	"os"

	"github.com/soner3/net-scan/dns/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dnsCmd represents the dns command
var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Perform DNS lookups for hosts listed in host file",
	Long: `The dns command performs DNS lookups for each host listed in the specified file.

It queries multiple DNS record types including:
  - A and AAAA (IPv4 and IPv6 addresses)
  - CNAME (Canonical name)
  - MX (Mail exchanger)
  - TXT (Text records)
  - NS (Name servers)

The results are printed in a structured format per host. Use the --file flag to specify the input file.

Example:
  net-scan dns --file hosts.txt

Each line in the input file should contain a single hostname.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := viper.GetString("file")
		return action.DnsAction(os.Stdout, filename)
	},
}

func init() {
	DnsCmd.SetErrPrefix("DNS Error:")
}
