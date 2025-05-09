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
	Short: "Performs DNS lookups for hosts in host file",
	Long: `The dns command queries DNS records for each host in a given file.

You can specify which record types to query using the --search flag.
Supported types include: cname, ip4, ip6, ns, mx, and txt.

Examples:
  net-scan dns --file targets.txt
  net-scan dns --file targets.txt --search cname,ip4,txt

Each result will include the requested DNS data, ordered by host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := viper.GetString("file")
		search := viper.GetStringSlice("dns.search")
		return action.DnsAction(os.Stdout, filename, &search)
	},
	SilenceUsage: true,
}

func init() {
	DnsCmd.SetErrPrefix("DNS Error:")

	DnsCmd.Flags().StringSliceP("search", "s", []string{"cname", "ip4", "ip6", "ns", "mx", "txt"}, "search for specific entries")

	viper.BindPFlag("dns.search", DnsCmd.Flags().Lookup("search"))
}
