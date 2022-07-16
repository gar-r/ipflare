package ipdetect

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// ChangeDetector runs in the background, sending notifications about
// changes in the public IP address of the current host.
type ChangeDetector struct {
	Frequency time.Duration
	C, Err    chan string
	previous  string
}

// NewChangeDetector initializes a ChangeDetector with the
// given polling frequency.
func NewChangeDetector(freq time.Duration) *ChangeDetector {
	return &ChangeDetector{
		Frequency: freq,
		C:         make(chan string),
		Err:       make(chan string),
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
				c.Err <- err.Error()
			} else if current != c.previous {
				c.previous = current
				c.C <- current
			}
		}
	}()
}

var ipLookupFunc = func() (string, error) {
	r := resty.New()
	res, err := r.NewRequest().
		SetQueryParam("format", "text").
		Get("https://api64.ipify.org")
	return res.String(), err
}
