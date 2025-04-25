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

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/scan"
)

type Config struct {
	filename      string
	ports         []int
	portRange     string
	network       string
	timeout       int
	resolvedPorts []int
}

func NewConfig(filename string, ports []int, portRange string, network string, timeout int) *Config {
	return &Config{
		filename:  filename,
		ports:     ports,
		portRange: portRange,
		network:   network,
		timeout:   timeout}
}

func ScanAction(out io.Writer, cfg *Config) error {
	hl := host.NewHostList()
	if err := validateConfig(cfg, hl); err != nil {
		return err
	}

	result := scan.Run(hl, &cfg.resolvedPorts, cfg.network, cfg.timeout)

	for _, res := range *result {
		output := fmt.Sprintf("%s:\n", res.Host)
		if res.NotFound {
			output += "\tNot Found\n\n"
		} else {
			portState := res.PortStates
			for _, ps := range *portState {
				output += fmt.Sprintf("\t%d/%s: %s\n", ps.Port, cfg.network, ps.Open)
			}
			output += "\n\n"
		}

		fmt.Fprint(out, output)

	}

	return nil
}

func validateConfig(cfg *Config, hl *host.HostList) error {
	if err := hl.Load(cfg.filename); err != nil {
		return err
	}
	if len(hl.Hosts) == 0 {
		return fmt.Errorf("host file is empty")
	}

	if len(cfg.ports) == 0 && cfg.portRange == "" {
		return fmt.Errorf("either --ports or --port-range must be set")
	}

	for _, p := range cfg.ports {
		if p < 1 {
			return fmt.Errorf("invalid port value: %d", p)
		}
	}

	var rangePorts []int
	if cfg.portRange != "" {
		var start, end int
		if _, err := fmt.Sscanf(cfg.portRange, "%d-%d", &start, &end); err != nil {
			return fmt.Errorf("invalid port range: %s", cfg.portRange)
		}
		if start < 1 || end < 1 || start >= end {
			return fmt.Errorf("invalid port range values: %s", cfg.portRange)
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

	cfg.resolvedPorts = append(rangePorts, cfg.ports...)
	sort.Ints(cfg.resolvedPorts)

	networks := []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

	if !slices.Contains(networks, cfg.network) {
		return fmt.Errorf("invalid network value: %s", cfg.network)
	}

	if cfg.timeout < 1 {
		return fmt.Errorf("timeout value must be greater than 0: %s", cfg.network)
	}
	return nil
}
