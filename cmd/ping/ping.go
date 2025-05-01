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
package ping

import (
	"os"
	"time"

	"github.com/soner3/net-scan/ping/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pingCmd represents the ping command
var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send ICMP ping requests to the hosts in the specified file",
	Long:  `The ping command sends ICMP echo requests to multiple hosts defined in a file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &action.Config{
			Filename:   viper.GetString("file"),
			Timeout:    viper.GetDuration("ping.timeout"),
			Interval:   viper.GetDuration("ping.interval"),
			Count:      viper.GetInt("ping.count"),
			Size:       viper.GetInt("ping.size"),
			Ttl:        viper.GetInt("ping.ttl"),
			Iface:      viper.GetString("ping.iface"),
			Tclass:     viper.GetInt("ping.tclass"),
			Priveleged: viper.GetBool("ping.privileged"),
		}
		return action.PingAction(os.Stdout, cfg)
	},
}

func init() {
	PingCmd.SetErrPrefix("Ping Error:\n\t")

	PingCmd.Flags().DurationP("timeout", "t", time.Second*10, "Timeout specifies a timeout before ping exits, regardless of how many packets have been received")
	PingCmd.Flags().DurationP("interval", "i", time.Second, "Interval is the wait time between each packet send")
	PingCmd.Flags().IntP("count", "c", 4, "Number to stop after sending (and receiving) Number echo packets. If this option is 0, it will operate until interrupted.")
	PingCmd.Flags().IntP("size", "s", 56, "Size of packet being sent")
	PingCmd.Flags().IntP("ttl", "l", 64, "Time to live (TTL)")
	PingCmd.Flags().StringP("iface", "I", "", "Network interface to use")
	PingCmd.Flags().IntP("tclass", "Q", -1, "The traffic class (type-of-service field for IPv4) value for future outgoing packets")
	PingCmd.Flags().Bool("privileged", false, "Use privileged raw socket")

	viper.BindPFlag("ping.timeout", PingCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("ping.interval", PingCmd.Flags().Lookup("interval"))
	viper.BindPFlag("ping.count", PingCmd.Flags().Lookup("count"))
	viper.BindPFlag("ping.size", PingCmd.Flags().Lookup("size"))
	viper.BindPFlag("ping.ttl", PingCmd.Flags().Lookup("ttl"))
	viper.BindPFlag("ping.iface", PingCmd.Flags().Lookup("iface"))
	viper.BindPFlag("ping.tclass", PingCmd.Flags().Lookup("tclass"))
	viper.BindPFlag("ping.privileged", PingCmd.Flags().Lookup("privileged"))
}
