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
	"os"

	"github.com/soner3/net-scan/host/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes one or more hosts using args: <host1> <host2> ... <hostN>",
	Long: `Deletes one or more hosts from the persistent host list.

You can provide hosts to delete directly as command-line arguments:
  net-scan host delete host1 host2 host3

Alternatively, you can provide a list of hosts via stdin:
  cat remove-these.txt | net-scan host delete

Each host must be on a separate line when passed via stdin.

The updated host list will be saved to the configured hosts file (default: net-scan.hosts).`,
	SilenceUsage: true,
	Aliases:      []string{"d"},
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := viper.GetString("file")
		return action.DeleteAction(os.Stdout, filename, args, os.Stdin)
	},
}

func init() {
	HostCmd.AddCommand(deleteCmd)
}
