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
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/soner3/net-scan/host"
	"github.com/soner3/net-scan/http"
)

type Config struct {
	Filename      string
	CallFrequency time.Duration
	Timeout       time.Duration
	Secure        bool
}

func NewConfig(filename string, callFrequency, timeout time.Duration, secure bool) *Config {
	return &Config{
		Filename:      filename,
		CallFrequency: callFrequency,
		Timeout:       timeout,
		Secure:        secure,
	}
}

var (
	ErrInvalidHTTP = errors.New("invalid HTTP config")
	ErrEmptyFile   = errors.New("host file is empty")
)

func (cfg *Config) validate() error {
	if cfg.Filename == "" {
		return fmt.Errorf("%w: filename must be set", ErrInvalidHTTP)
	}

	if _, err := os.Stat(cfg.Filename); err != nil {
		return fmt.Errorf("%w: file not found: %s", ErrInvalidHTTP, err.Error())
	}

	hl := host.NewHostList()
	if err := hl.Load(cfg.Filename); err != nil {
		return err
	}
	if len(hl.Hosts) == 0 {
		return fmt.Errorf("%w: %s", ErrEmptyFile, cfg.Filename)
	}

	if cfg.CallFrequency <= 0 {
		return fmt.Errorf("%w: call-frequency must be > 0", ErrInvalidHTTP)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("%w: timeout must be > 0", ErrInvalidHTTP)
	}

	return nil
}

func HttpAction(out io.Writer, cfg *Config) error {
	if err := cfg.validate(); err != nil {
		return err
	}

	hl := host.NewHostList()
	if err := hl.Load(cfg.Filename); err != nil {
		return err
	}
	for _, h := range hl.Hosts {
		if err := http.Run(out, h, cfg.Secure, cfg.CallFrequency, cfg.Timeout); err != nil {
			return err
		}
	}
	return nil
}
