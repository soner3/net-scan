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
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/soner3/net-scan/host"
)

func setup(t *testing.T, name string) string {
	tf, err := os.CreateTemp(".", fmt.Sprintf("%s*", name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		tf.Close()
		os.Remove(tf.Name())
	})
	return tf.Name()
}

func TestScanActionValidation(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         *Config
		expectedErr error
		expectedOut *[]int
	}{
		{"ValidateFileNotFound", NewConfig("not-found.txt", []int{}, "", "", 1, ""), os.ErrNotExist, nil},
		{"ValidateHostFileEmpty", NewConfig("", []int{}, "", "", 1, ""), ErrEmpty, nil},
		{"ValidatePortsAndRangeEmpty", NewConfig("", []int{}, "", "", 1, ""), ErrEmpty, nil},
		{"ValidatePorts", NewConfig("", []int{1, -2}, "", "", 1, ""), ErrValue, nil},
		{"ValidatePortRangeFormat", NewConfig("", []int{}, "78655", "", 1, ""), ErrFormat, nil},
		{"ValidatePortRangeValue", NewConfig("", []int{}, "-10-23", "", 1, ""), ErrValue, nil},
		{"ValidateNetwork", NewConfig("", []int{1}, "", "khu", 1, ""), ErrValue, nil},
		{"ValidateTimeout", NewConfig("", []int{}, "10-23", "tcp", -1, ""), ErrValue, nil},
		{"ValidateFilterErr", NewConfig("", []int{1}, "", "tcp", 1, "dfgb"), ErrValue, nil},
		{"ValidateSuccessWithoutFilter", NewConfig("", []int{1}, "", "tcp", 1, ""), nil, &[]int{1}},
		{"ValidateSuccessWithoutOpen", NewConfig("", []int{1}, "", "tcp", 1, "open"), nil, &[]int{1}},
		{"ValidateSuccessWithoutTimeout", NewConfig("", []int{1}, "", "tcp", 1, "timeout"), nil, &[]int{1}},
		{"ValidateSuccessWithoutClosed", NewConfig("", []int{1}, "", "tcp", 1, "closed"), nil, &[]int{1}},
		{"ValidateSuccessUniquePorts", NewConfig("", []int{1, 2, 3, 4, 5}, "1-5", "tcp", 1, ""), nil, &[]int{1, 2, 3, 4, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name != "ValidateFileNotFound" {
				tc.cfg.filename = setup(t, tc.name)
			} else {
				defer os.Remove(tc.cfg.filename)
			}

			if tc.name != "ValidateHostFileEmpty" && tc.name != "ValidateFileNotFound" {
				hl := host.NewHostList()
				hl.Add("host1")
				hl.Save(tc.cfg.filename)
			}

			out, outErr := tc.cfg.validate()

			if tc.expectedErr != nil && tc.expectedOut == nil {
				if !errors.Is(outErr, tc.expectedErr) {
					t.Errorf("Expected %q, got %q instead", tc.expectedErr, outErr)
				}

				if out != tc.expectedOut {
					t.Errorf("Expected %v, got %v instead", tc.expectedOut, out)
				}

			} else if tc.expectedOut != nil && tc.expectedErr == nil {
				if outErr != tc.expectedErr {
					t.Errorf("Expected %v, got %v instead", tc.expectedErr, outErr)
				}

				if len(*out) != len(*tc.expectedOut) {
					t.Errorf("Expected %d, got %d instead", len(*tc.expectedOut), len(*out))
				}

				for i, expOutPort := range *tc.expectedOut {
					if expOutPort != (*out)[i] {
						t.Errorf("Expected %d, got %d instead", expOutPort, (*out)[i])
					}
				}
			} else {
				t.Errorf("Unexpected behavior")
			}

		})
	}

}

func TestScanAction(t *testing.T) {
	ports := []int{}
	closedPort := 0
	for i := range 2 {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		_, portString, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portString)
		if err != nil {
			t.Fatal(err)
		}
		ports = append(ports, port)
		if i == 1 {
			closedPort = port
			ln.Close()
		}
	}

	tc := []struct {
		host  string
		ports []int
		found bool
	}{
		{"192.0.2.1", ports, true},
		{"localhost", ports, true},
		{"unknown", nil, false},
	}
	slices.Sort(ports)
	cfg := NewConfig("", ports, "", "tcp", time.Second, "")
	cfg.filename = setup(t, "ScanActionTest")

	hl := host.NewHostList()
	for _, c := range tc {
		hl.Add(c.host)
	}
	hl.Save(cfg.filename)
	expectedOut := ""
	for _, c := range tc {
		if !c.found {
			expectedOut += fmt.Sprintf("%s:\n\tNot Found\n\n", c.host)
			continue
		}
		expectedOut += fmt.Sprintf("%s:\n", c.host)

		if c.host == "192.0.2.1" {
			for _, p := range c.ports {
				expectedOut += fmt.Sprintf("\t%d/tcp: timeout\n", p)
			}
		} else {
			for _, p := range c.ports {
				if p == closedPort {
					expectedOut += fmt.Sprintf("\t%d/tcp: closed\n", p)
				} else {
					expectedOut += fmt.Sprintf("\t%d/tcp: open\n", p)

				}
			}
		}
		expectedOut += "\n"

	}

	var out bytes.Buffer
	err := ScanAction(&out, cfg)
	if err != nil {
		t.Errorf("Expected nil, got %q instead", err)
	}

	if out.String() != expectedOut {
		t.Errorf("Expected %s, got %s instead", expectedOut, out.String())

	}

}
