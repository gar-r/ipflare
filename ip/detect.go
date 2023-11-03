package ip

import (
	"time"
)

var DefaultIpProvider IpProvider = &Ipify{}

// ChangeDetector runs in the background, sending notifications about
// changes in the public IP address of the current host.
type ChangeDetector struct {
	Provider  IpProvider
	C         chan string
	Err       chan error
	previous  string
	Frequency time.Duration
}

// IpProvider provides the public IP address.
type IpProvider interface {
	GetPublicIp() (string, error)
}

// NewChangeDetector initializes a ChangeDetector with the
// given polling frequency and the default ip provider.
func NewChangeDetector(freq time.Duration) *ChangeDetector {
	return &ChangeDetector{
		Provider:  DefaultIpProvider,
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
			current, err := c.Provider.GetPublicIp()
			if err != nil {
				c.Err <- err
			} else if current != c.previous {
				c.previous = current
				c.C <- current
			}
		}
	}()
}
