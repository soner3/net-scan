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
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var HostCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage predefined hosts for quick access",
	Long: `Manage a list of predefined hosts used across different net-scan commands.

This command allows you to add, delete, or list commonly used targets
so you don't have to retype IPs or domains for every scan.

Subcommands:
  - list:    Show all saved hosts
  - add:     Add a new host to the list
  - delete:  Remove a host from the list

Example usage:
  net-scan host add myserver 192.168.1.10
  net-scan host list
  net-scan host delete myserver`,
}

func init() {
	HostCmd.SetErrPrefix("Host Error:\n\t")
}
