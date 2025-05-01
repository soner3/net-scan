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
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func Run(out io.Writer, host string, secure bool, callFrequency, timeout time.Duration) error {
	var url string
	if secure {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = fmt.Sprintf("http://%s", host)
	}

	httpCaller := probing.NewHttpCaller(
		url,
		probing.WithHTTPCallerCallFrequency(callFrequency),
		probing.WithHTTPCallerOnResp(func(suite *probing.TraceSuite, info *probing.HTTPCallInfo) {
			fmt.Fprintf(out, "\tgot resp, status code: %d, latency: %s\n",
				info.StatusCode,
				suite.GetGeneralEnd().Sub(suite.GetGeneralStart()),
			)
		}),
		probing.WithHTTPCallerTimeout(timeout),
	)

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		httpCaller.Stop()
	}()

	fmt.Fprintf(out, "%s:\n", url)

	if _, err := net.LookupHost(host); err != nil {
		fmt.Fprint(out, "\tNot Found\n\n")
		return nil
	}

	httpCaller.Run()
	fmt.Fprintf(out, "\n")
	return nil
}
