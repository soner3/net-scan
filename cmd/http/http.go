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
package http

import (
	"os"
	"time"

	"github.com/soner3/net-scan/http/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// httpCmd represents the http command
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Perform HTTP availability checks on hosts from file",
	Long: `The http command sends periodic HTTP requests to hosts defined in a file.

You can configure the frequency and timeout of the requests. Optionally, you can
enable HTTPS with the --secure flag.

Example:
  net-scan http --call-frequency 1s --timeout 5s --secure
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &action.Config{
			Filename:      viper.GetString("file"),
			CallFrequency: viper.GetDuration("http.call-frequency"),
			Timeout:       viper.GetDuration("http.timeout"),
			Secure:        viper.GetBool("http.secure"),
		}
		return action.HttpAction(os.Stdout, cfg)
	},
}

func init() {
	HttpCmd.SetErrPrefix("HTTP Error:")

	HttpCmd.Flags().DurationP("call-frequency", "r", 1*time.Second, "Time between HTTP requests per host")
	HttpCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Request timeout duration")
	HttpCmd.Flags().BoolP("secure", "s", true, "Use HTTPS instead of HTTP")

	viper.BindPFlag("http.call-frequency", HttpCmd.Flags().Lookup("call-frequency"))
	viper.BindPFlag("http.timeout", HttpCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("http.secure", HttpCmd.Flags().Lookup("secure"))
}
