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
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Config struct {
	Count      int
	Size       int
	Interval   time.Duration
	Timeout    time.Duration
	TTL        int
	Iface      string
	Privileged bool
	TClass     int
}

func NewConfig(count, size int, interval, timeout time.Duration, ttl int, iface string, privileged bool, tclass int) *Config {
	return &Config{
		Count:      count,
		Size:       size,
		Interval:   interval,
		Timeout:    timeout,
		TTL:        ttl,
		Iface:      iface,
		Privileged: privileged,
		TClass:     tclass,
	}
}

func Run(out io.Writer, host string, cfg *Config) error {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		if _, err := net.LookupHost(host); err != nil {
			fmt.Fprintf(out, "%s:\n\tNot Found\n\n", host)
			return nil
		} else {
			return err
		}
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	pinger.Count = cfg.Count
	pinger.Size = cfg.Size
	pinger.Interval = cfg.Interval
	pinger.Timeout = cfg.Timeout
	pinger.TTL = cfg.TTL
	pinger.InterfaceName = cfg.Iface
	pinger.SetPrivileged(cfg.Privileged)
	pinger.SetTrafficClass(uint8(cfg.TClass))

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Fprintf(out, "\t%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Fprintf(out, "\t%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Fprintf(out, "\n\t--- %s ping statistics ---\n", stats.Addr)
		fmt.Fprintf(out, "\t%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Fprintf(out, "\tround-trip min/avg/max/stddev = %v/%v/%v/%v\n\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Fprintf(out, "PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		return fmt.Errorf("failed to ping target host: %w", err)
	}
	return nil
}
