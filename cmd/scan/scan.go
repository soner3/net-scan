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
	"os"

	"github.com/soner3/net-scan/scan/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan ports on hosts defined in a host file",
	Long: `The scan command connects to specified ports on target hosts
using a selected network protocol (e.g., tcp, udp, etc.). 
You can define specific ports or port ranges and apply filters to show only open, closed, or timeout states.

Supported network protocols include: 
tcp, tcp4, tcp6, udp, udp4, udp6, ip, ip4, ip6, unix, unixgram, unixpacket.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := viper.GetString("file")
		ports := viper.GetIntSlice("scan.ports")
		portRange := viper.GetString("scan.port-range")
		network := viper.GetString("scan.network")
		timeout := viper.GetInt("scan.timeout")
		filter := viper.GetString("scan.filter-state")

		return action.ScanAction(os.Stdout, action.NewConfig(filename, ports, portRange, network, timeout, filter))
	},
}

func init() {
	ScanCmd.SetErrPrefix("Scan Error:\n\t")

	ScanCmd.Flags().IntSliceP("ports", "p", []int{}, "Ports to scan on the target hosts (e.g., 22,80,443)")
	ScanCmd.Flags().StringP("port-range", "r", "", "Port range to scan on the target hosts (e.g., 20-100)")
	ScanCmd.Flags().StringP("network", "n", "tcp", "Network protocol to use (tcp, udp, tcp4, tcp6, udp4, udp6, ip, ip4, ip6, unix, unixgram, unixpacket)")
	ScanCmd.Flags().IntP("timeout", "t", 1000, "Timeout per port in milliseconds")
	ScanCmd.Flags().StringP("filter-state", "s", "", "Filter scanned results by port state (open, closed, timeout)")

	viper.BindPFlag("scan.ports", ScanCmd.Flags().Lookup("ports"))
	viper.BindPFlag("scan.port-range", ScanCmd.Flags().Lookup("port-range"))
	viper.BindPFlag("scan.network", ScanCmd.Flags().Lookup("network"))
	viper.BindPFlag("scan.timeout", ScanCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("scan.filter-state", ScanCmd.Flags().Lookup("filter-state"))
}
