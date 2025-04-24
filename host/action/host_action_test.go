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
package action_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/host/action"
)

func setup(t *testing.T, hosts []string, initList bool) string {
	tf, err := os.CreateTemp(".", "test-add-and-delete-action*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tf.Close()
		if err := os.Remove(tf.Name()); err != nil {
			t.Fatal(err)
		}
	})
	hl := host.NewHostList()

	if initList {

		for _, h := range hosts {
			if err := hl.Add(h); err != nil {
				t.Fatal(err)
			}
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}

	}

	return tf.Name()
}

func TestHostAddAndDeleteAction(t *testing.T) {
	testData := []struct {
		name        string
		args        []string
		expectedOut string
		initList    bool
		execAction  func(out io.Writer, filename string, args []string, input io.Reader) error
	}{
		{"AddAction", []string{"host1", "host2", "host3"}, "Added host: host1\nAdded host: host2\nAdded host: host3\n", false, action.AddAction},
		{"DeleteAction", []string{"host1", "host2", "host3"}, "Deleted host: host1\nDeleted host: host2\nDeleted host: host3\n", true, action.DeleteAction},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			var out bytes.Buffer
			filename := setup(t, td.args, td.initList)

			err := td.execAction(&out, filename, td.args, os.Stdin)
			if err != nil {
				t.Errorf("Expected no error, got %q instead", err)
			}

			if out.String() != td.expectedOut {
				t.Errorf("Expected %s, got %s instead", td.expectedOut, out.String())
			}

		})
	}
}

func TestHostListAction(t *testing.T) {
	args := []string{"host1", "host2", "host3"}
	expectedOut := "host1\nhost2\nhost3\n"

	filename := setup(t, args, true)
	var out bytes.Buffer

	if err := action.ListAction(&out, filename); err != nil {
		t.Errorf("Expected no error, got %q instead", err)
	}

	if out.String() != expectedOut {
		t.Errorf("Expected %s, got %s instead", expectedOut, out.String())
	}

}

func TestIntegration(t *testing.T) {
	var out bytes.Buffer

	hosts := []string{
		"host1",
		"host2",
		"host3",
	}
	delHosts := []string{
		"host1",
		"host3",
	}
	hostsEnd := []string{
		"host2",
	}
	filename := setup(t, hosts, false)

	expectedOut := ""
	expectedOut += makeActionOutput("Added host: %s\n", hosts)
	expectedOut += makeActionOutput("%s\n", hosts)
	expectedOut += makeActionOutput("Deleted host: %s\n", delHosts)
	expectedOut += makeActionOutput("%s\n", hostsEnd)

	if err := action.AddAction(&out, filename, hosts, os.Stdin); err != nil {
		t.Errorf("Expected no error, got %q instead", err)
	}

	if err := action.ListAction(&out, filename); err != nil {
		t.Errorf("Expected no error, got %q instead", err)
	}

	if err := action.DeleteAction(&out, filename, delHosts, os.Stdin); err != nil {
		t.Errorf("Expected no error, got %q instead", err)
	}

	if err := action.ListAction(&out, filename); err != nil {
		t.Errorf("Expected no error, got %q instead", err)
	}

	if out.String() != expectedOut {
		t.Errorf("Expected %s, got %s instead", expectedOut, out.String())
	}

}

func makeActionOutput(format string, hosts []string) string {
	out := ""
	for _, h := range hosts {
		out += fmt.Sprintf(format, h)
	}
	return out
}

func TestHostStdinActions(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", t.Name(), "./../../")

	tf, err := os.CreateTemp(".", "test-host-stdin*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tf.Close()
		if err := os.Remove(tf.Name()); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(t.Name()); err != nil {
			t.Fatal(err)
		}
	})

	testData := []struct {
		name        string
		args        []string
		input       string
		expectedOut string
		cmd         string
	}{
		{"StdinAddAction", []string{"host1", "host2", "host3"}, "host4\nhost5\n", "Added host: host4\nAdded host: host5\nAdded host: host1\nAdded host: host2\nAdded host: host3\n", "add"},
		{"StdinDeleteAction", []string{"host1", "host2", "host3"}, "host4\nhost5\n", "Deleted host: host4\nDeleted host: host5\nDeleted host: host1\nDeleted host: host2\nDeleted host: host3\n", "delete"},
	}
	if err := buildCmd.Run(); err != nil {
		t.Fatal(err)
	}

	for _, td := range testData {
		args := append([]string{"-f", tf.Name(), "host", td.cmd}, td.args...)
		cmd := exec.Command("./"+t.Name(), args...)
		cmd.Stdin = strings.NewReader(td.input)

		out, err := cmd.Output()
		if err != nil {
			t.Errorf("Expected no error, got %q instead", err)
		}

		if string(out) != td.expectedOut {
			t.Errorf("Expected %s, got %s instead", td.expectedOut, string(out))
		}
	}
}
