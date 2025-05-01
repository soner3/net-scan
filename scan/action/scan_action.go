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
	"slices"
	"time"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/scan"
	"github.com/soner3/net-scan/util"
)

var (
	ErrEmpty  = errors.New("empty value")
	ErrValue  = errors.New("invalid value")
	ErrFormat = errors.New("invalid value format")
)

var networks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

var cmp = func(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

type Config struct {
	filename  string
	ports     []int
	portRange string
	network   string
	timeout   time.Duration
	filter    string
}

func NewConfig(filename string, ports []int, portRange string, network string, timeout time.Duration, filter string) *Config {
	return &Config{
		filename:  filename,
		ports:     ports,
		portRange: portRange,
		network:   network,
		timeout:   timeout,
		filter:    filter,
	}
}
func (cfg *Config) validate() (*[]int, error) {
	if cfg.filename == "" {
		return nil, fmt.Errorf("%w: filename must be provided", ErrEmpty)
	}

	file, err := os.Open(cfg.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hl := host.NewHostList()
	if err := hl.Load(cfg.filename); err != nil {
		return nil, err
	}
	if len(hl.Hosts) < 1 {
		return nil, fmt.Errorf("%w: host file is empty", ErrEmpty)
	}

	if len(cfg.ports) == 0 && cfg.portRange == "" {
		return nil, fmt.Errorf("%w: either --ports or --port-range must be set", ErrEmpty)
	}

	rangePorts := util.Set[int]{}
	for _, p := range cfg.ports {
		if p < 1 || p > 65535 {
			return nil, fmt.Errorf("%w: port %d is out of valid range (1–65535)", ErrValue, p)
		}
		rangePorts.Add(p)
	}

	if cfg.portRange != "" {
		var start, end int
		if _, err := fmt.Sscanf(cfg.portRange, "%d-%d", &start, &end); err != nil {
			return nil, fmt.Errorf("%w: port-range format must be start-end", ErrFormat)
		}
		if start < 1 || end < 1 || start >= end || end > 65535 {
			return nil, fmt.Errorf("%w: invalid port range %s", ErrValue, cfg.portRange)
		}
		for p := start; p <= end; p++ {
			rangePorts.Add(p)
		}
	}

	if !slices.Contains(networks, cfg.network) {
		return nil, fmt.Errorf("%w: unsupported network '%s'", ErrValue, cfg.network)
	}

	if cfg.timeout <= 0 {
		return nil, fmt.Errorf("%w: timeout must be greater than 0", ErrValue)
	}

	validFilters := []string{"open", "closed", "timeout", ""}
	if !slices.Contains(validFilters, cfg.filter) {
		return nil, fmt.Errorf("%w: unknown filter '%s'", ErrValue, cfg.filter)
	}

	return rangePorts.ToSortedSlice(cmp), nil
}

func ScanAction(out io.Writer, cfg *Config) error {
	resolvedPorts, err := cfg.validate()
	if err != nil {
		return err
	}

	hl := host.NewHostList()
	if err := hl.Load(cfg.filename); err != nil {
		return err
	}

	result := scan.Run(hl, resolvedPorts, cfg.network, cfg.timeout)

	for _, res := range *result {
		output := fmt.Sprintf("%s:\n", res.Host)
		if res.NotFound {
			output += "\tNot Found\n"
		} else {
			portState := res.PortStates
			for _, ps := range *portState {
				if cfg.filter != "" {
					if cfg.filter == ps.Open.String() {
						output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, &ps.Open)
					}
				} else {
					output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, &ps.Open)
				}

			}
		}
		output += "\n"

		_, err := fmt.Fprint(out, output)
		if err != nil {
			return err
		}

	}

	return nil
}
