package ipdetect

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// ChangeDetector runs in the background, sending notifications about
// changes in the public IP address of the current host.
type ChangeDetector struct {
	Frequency time.Duration
	C         chan string
	Err       chan error
	previous  string
}

// NewChangeDetector initializes a ChangeDetector with the
// given polling frequency.
func NewChangeDetector(freq time.Duration) *ChangeDetector {
	return &ChangeDetector{
		Frequency: freq,
		C:         make(chan string),
		Err:       make(chan error),
	}
}

// Start the change detection goroutine.
// Changes to the public IP will be sent to channel C
// Errors will be sent to channel Err
func (c *ChangeDetector) Start() {
	t := time.NewTicker(c.Frequency)
	go func() {
		for range t.C {
			current, err := ipLookupFunc()
			if err != nil {
				c.Err <- err
			} else if current != c.previous {
				c.previous = current
				c.C <- current
			}
		}
	}()
}

var ipLookupFunc = func() (string, error) {
	r := resty.New()
	r.SetTransport(&http.Transport{
		TLSHandshakeTimeout: 30 * time.Second,
	})
	res, err := r.NewRequest().
		SetQueryParam("format", "text").
		Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	var ip = res.String()
	if net.ParseIP(ip) == nil {
		return "", errors.New("response is not a valid ip")
	}
	return ip, nil
}
