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
package host_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/soner3/net-scan/host"
)

func TestAdd(t *testing.T) {
	testData := []struct {
		name        string
		host        string
		expectedErr error
		expectedLen int
	}{
		{"AddNew", "host2", nil, 2},
		{"AddExisting", "host1", host.ErrExists, 1},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			hl := host.NewHostList()

			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			out := hl.Add(td.host)

			if td.expectedErr != nil {
				if !errors.Is(out, host.ErrExists) {
					t.Errorf("Expected %q, got %q instead", td.expectedErr, out)
				}

				if len(hl.Hosts) != td.expectedLen {
					t.Errorf("Expected %d got %d instead", td.expectedLen, len(hl.Hosts))
				}

			}

			if len(hl.Hosts) != td.expectedLen {
				t.Errorf("Expected %d got %d instead", td.expectedLen, len(hl.Hosts))
			}
		})

	}

}

func TestRemove(t *testing.T) {
	testData := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{"RemoveExisting", "host1", 0, nil},
		{"RemoveNotExisting", "host2", 1, host.ErrNotExists},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			hl := host.NewHostList()

			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			out := hl.Remove(td.host)

			if td.expectedErr != nil {
				if !errors.Is(out, host.ErrNotExists) {
					t.Errorf("Expected %q, got %q instead", td.expectedErr, out)
				}

				if len(hl.Hosts) != td.expectedLen {
					t.Errorf("Expected %d got %d instead", td.expectedLen, len(hl.Hosts))
				}

			}

			if len(hl.Hosts) != td.expectedLen {
				t.Errorf("Expected %d got %d instead", td.expectedLen, len(hl.Hosts))
			}

		})
	}

}

func TestString(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	hl := host.NewHostList()
	for _, h := range hosts {
		if err := hl.Add(h); err != nil {
			t.Fatal(err)
		}
	}

	expectedOut := "host1\nhost2\nhost3\n"

	if fmt.Sprint(hl) != expectedOut {
		t.Errorf("Expected %s, got %s instead", expectedOut, fmt.Sprint(hl))
	}

}

func TestLoadSave(t *testing.T) {
	tf, err := os.CreateTemp(".", "load-save-test*.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tf.Close()
		err := os.Remove(tf.Name())
		if err != nil {
			t.Fatal(err)
		}
	})

	hosts := []string{"host1", "host2", "host3"}
	hl := host.NewHostList()
	for _, h := range hosts {
		if err := hl.Add(h); err != nil {
			t.Fatal(err)
		}
	}

	err = hl.Save(tf.Name())
	if err != nil {
		t.Fatal(err)
	}

	hl2 := host.NewHostList()
	err = hl2.Load(tf.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(hl2.Hosts) != len(hl.Hosts) {
		t.Errorf("Expected %d, got %d instead", len(hl.Hosts), len(hl2.Hosts))
	}

	expectedOut := "host1\nhost2\nhost3\n"

	if fmt.Sprint(hl2) != expectedOut {
		t.Errorf("Expected %s, got %s instead", expectedOut, fmt.Sprint(hl2))
	}
}

func TestLoadNoFile(t *testing.T) {
	hl := host.NewHostList()
	err := hl.Load("test-load-no-file")

	if err != nil {
		t.Errorf("Expected nil, got %q instead", err)
	}

	if len(hl.Hosts) != 0 {
		t.Errorf("Expected 0, got %d instead", len(hl.Hosts))
	}

}
