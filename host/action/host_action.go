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

	"github.com/soner3/net-scan/host"
)

func AddAction(out io.Writer, filename string, args []string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}
		_, err := fmt.Fprintln(out, "Added host:", h)
		if err != nil {
			return err
		}
	}

	return hl.Save(filename)
}

func DeleteAction(out io.Writer, filename string, args []string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}
		_, err := fmt.Fprintln(out, "Deleted host:", h)
		if err != nil {
			return err
		}
	}

	return hl.Save(filename)
}

func ListAction(out io.Writer, filename string) error {
	hl := host.NewHostList()
	if err := hl.Load(filename); err != nil {
		return err
	}
	_, err := fmt.Fprint(out, hl)
	return err
}
