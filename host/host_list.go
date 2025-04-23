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
package host

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host does not exist in the list")
)

type HostList struct {
	Hosts []string
}

func NewHostList() *HostList {
	return &HostList{}
}

// Add host to host list
func (hl *HostList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	hl.Hosts = append(hl.Hosts, host)

	return nil
}

// Remove host from host list
func (hl *HostList) Remove(host string) error {
	found, i := hl.search(host)
	if !found {
		return fmt.Errorf("%w: %s", ErrNotExists, host)
	}
	hl.Hosts = slices.Delete(hl.Hosts, i, i+1)
	return nil
}

// Load hosts from file
func (hl *HostList) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err

	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		host := scanner.Text()
		hl.Add(host)
	}

	return nil
}

// Save hosts to file
func (hl *HostList) Save(file string) error {
	output := fmt.Sprint(hl)
	return os.WriteFile(file, []byte(output), 0644)
}

// Search host in host list
func (hl *HostList) search(host string) (bool, int) {
	slices.Sort(hl.Hosts)
	i := slices.Index(hl.Hosts, host)
	if i != -1 && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1

}

// String method
func (hl *HostList) String() string {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	return output
}
