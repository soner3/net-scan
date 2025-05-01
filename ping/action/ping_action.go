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
	"errors"
	"fmt"
	"io"
	"os"
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

var (
	ErrEmptyFile   = errors.New("host file is empty")
	ErrInvalidPing = errors.New("invalid ping config")
)

func (cfg *Config) validate() error {
	if cfg.Filename == "" {
		return fmt.Errorf("%w: filename must be set", ErrInvalidPing)
	}

	if _, err := os.Stat(cfg.Filename); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidPing, err.Error())
	}

	hl := host.NewHostList()
	if err := hl.Load(cfg.Filename); err != nil {
		return err
	}
	if len(hl.Hosts) == 0 {
		return ErrEmptyFile
	}

	if cfg.Count < 0 {
		return fmt.Errorf("%w: count must be ≥ 0", ErrInvalidPing)
	}
	if cfg.Size < 0 || cfg.Size > 65500 {
		return fmt.Errorf("%w: packet size must be between 0 and 65500", ErrInvalidPing)
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("%w: interval must be > 0", ErrInvalidPing)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("%w: timeout must be > 0", ErrInvalidPing)
	}
	if cfg.Ttl <= 0 || cfg.Ttl > 255 {
		return fmt.Errorf("%w: ttl must be in range 1–255", ErrInvalidPing)
	}
	if cfg.Tclass < -1 || cfg.Tclass > 255 {
		return fmt.Errorf("%w: tclass must be between -1 and 255", ErrInvalidPing)
	}

	return nil
}

func PingAction(out io.Writer, cfg *Config) error {
	if err := cfg.validate(); err != nil {
		return err
	}

	hl := host.NewHostList()
	if err := hl.Load(cfg.Filename); err != nil {
		return err
	}
	for _, h := range hl.Hosts {
		pingCfg := ping.NewConfig(
			cfg.Count,
			cfg.Size,
			cfg.Interval,
			cfg.Timeout,
			cfg.Ttl,
			cfg.Iface,
			cfg.Priveleged,
			cfg.Tclass,
		)
		if err := ping.Run(out, h, pingCfg); err != nil {
			return err
		}
	}
	return nil
}
