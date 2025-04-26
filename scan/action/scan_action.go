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
	"sort"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/scan"
)

var (
	ErrEmpty  = errors.New("flag not set")
	ErrValue  = errors.New("invalid flag value")
	ErrFormat = errors.New("invalid flag value format")
)

var networks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

type Config struct {
	filename  string
	ports     []int
	portRange string
	network   string
	timeout   int
	filter    string
}

func NewConfig(filename string, ports []int, portRange string, network string, timeout int, filter string) *Config {
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
	f, err := os.Open(cfg.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if len(cfg.ports) == 0 && cfg.portRange == "" {
		return nil, fmt.Errorf("%q: either --ports or --port-range must be set", ErrEmpty)
	}

	if cfg.filter != "closed" && cfg.filter != "open" && cfg.filter != "timeout" && cfg.filter != "" {
		return nil, fmt.Errorf("%q: %s", ErrValue, cfg.filter)
	}

	if !slices.Contains(networks, cfg.network) {
		return nil, fmt.Errorf("%q: %s", ErrValue, cfg.network)
	}

	if cfg.timeout < 1 {
		return nil, fmt.Errorf("%q: %s", ErrValue, cfg.network)
	}

	for _, p := range cfg.ports {
		if p < 1 {
			return nil, fmt.Errorf("%q: %d", ErrValue, p)
		}
	}

	var rangePorts []int

	if cfg.portRange != "" {
		var start, end int
		if _, err := fmt.Sscanf(cfg.portRange, "%d-%d", &start, &end); err != nil {
			return nil, fmt.Errorf("%q: %s", ErrFormat, cfg.portRange)
		}
		if start < 1 || end < 1 || start >= end {
			return nil, fmt.Errorf("%q: %s", ErrValue, cfg.portRange)
		}

		rangePorts = make([]int, 0, end-start+1)
		for p := start; p <= end; p++ {
			for i, cfgP := range cfg.ports {
				if cfgP == p {
					cfg.ports = slices.Delete(cfg.ports, i, i+1)
				}
			}
			rangePorts = append(rangePorts, p)
		}
	}

	rangePorts = append(rangePorts, cfg.ports...)
	sort.Ints(rangePorts)

	return &rangePorts, nil
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
				switch true {
				case cfg.filter == "timeout":
					if ps.Timeout {
						output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, ps)
					}
				case cfg.filter == "closed":
					if !ps.Open && !ps.Timeout {
						output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, ps)
					}
				case cfg.filter == "open":
					if ps.Open {
						output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, ps)
					}
				default:
					output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, ps)
				}

			}
			output += "\n"
		}

		fmt.Fprint(out, output)

	}

	return nil
}
