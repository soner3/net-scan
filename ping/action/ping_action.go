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
	"io"
	"time"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/ping"
)

type Config struct {
	Filename   string
	Timeout    time.Duration
	Interval   time.Duration
	Count      int
	Size       int
	Ttl        int
	Iface      string
	Tclass     int
	Priveleged bool
}

func NewConfig(filename string, timeout, interval time.Duration, count, size, ttl int, iface string, tclass int, privileged bool) *Config {
	return &Config{
		Filename:   filename,
		Timeout:    timeout,
		Interval:   interval,
		Count:      count,
		Size:       size,
		Ttl:        ttl,
		Iface:      iface,
		Tclass:     tclass,
		Priveleged: privileged,
	}
}

func PingAction(out io.Writer, cfg *Config) error {
	hl := host.NewHostList()
	if err := hl.Load(cfg.Filename); err != nil {
		return err
	}
	for _, h := range hl.Hosts {
		pingCfg := &ping.Config{
			Count:      cfg.Count,
			Size:       cfg.Size,
			Interval:   cfg.Interval,
			Timeout:    cfg.Timeout,
			TTL:        cfg.Ttl,
			Iface:      cfg.Iface,
			Privileged: cfg.Priveleged,
			TClass:     cfg.Tclass,
		}
		if err := ping.Run(out, h, pingCfg); err != nil {
			return err
		}
	}
	return nil
}
