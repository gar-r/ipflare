package ipdetect

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewChangeDetector(t *testing.T) {
	t.Run("frequency initialized", func(t *testing.T) {
		freq := 5 * time.Second
		cd := NewChangeDetector(freq)
		assert.Equal(t, freq, cd.Frequency)
	})

	t.Run("channels initialized", func(t *testing.T) {
		cd := NewChangeDetector(time.Second)
		assert.NotNil(t, cd.C)
		assert.NotNil(t, cd.Err)
	})
}

func TestStart(t *testing.T) {

	t.Run("error emitted", func(t *testing.T) {
		ipLookupFunc = func() (string, error) {
			return "", errors.New("test error")
		}
		cd := NewChangeDetector(time.Millisecond)
		cd.Start()
		err := <-cd.Err
		assert.EqualError(t, err, "test error")
	})

	t.Run("change emitted", func(t *testing.T) {
		ipLookupFunc = func() (string, error) {
			return "test2", nil
		}
		cd := NewChangeDetector(time.Millisecond)
		cd.previous = "test1"
		cd.Start()
		ip := <-cd.C
		assert.Equal(t, "test2", ip)
	})

	t.Run("nothing emitted when no change", func(t *testing.T) {
		val := "test"
		ipLookupFunc = func() (string, error) {
			return val, nil
		}
		cd := NewChangeDetector(time.Millisecond)
		cd.previous = val
		timeout := time.After(10 * time.Millisecond)
		cd.Start()
		for {
			select {
			case <-cd.C:
				assert.FailNow(t, "unexpected event on channel")
			case <-timeout:
				return
			}
		}
	})

}
