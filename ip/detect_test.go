package ip

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	t.Run("provider initialized", func(t *testing.T) {
		cd := NewChangeDetector(time.Second)
		assert.NotNil(t, cd.Provider)
	})
}

func TestStart(t *testing.T) {
	cd := NewChangeDetector(time.Millisecond)

	t.Run("error emitted", func(t *testing.T) {
		p := &TestProvider{}
		p.On("GetPublicIp").Return("", errors.New("test error"))
		cd.Provider = p
		cd.Start()
		err := <-cd.Err
		assert.EqualError(t, err, "test error")
	})

	t.Run("change emitted", func(t *testing.T) {
		p := &TestProvider{}
		p.On("GetPublicIp").Return("127.0.0.1", nil)
		cd.Provider = p
		cd.Start()
		ip := <-cd.C
		assert.Equal(t, "127.0.0.1", ip)
	})

	t.Run("nothing emitted when no change", func(t *testing.T) {
		ip := "127.0.0.1"
		p := &TestProvider{}
		p.On("GetPublicIp").Return(ip, nil)
		cd.previous = ip
		cd.Provider = p
		cd.Start()
		timeout := time.After(10 * time.Millisecond)
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

type TestProvider struct {
	mock.Mock
}

func (t *TestProvider) GetPublicIp() (string, error) {
	args := t.Called()
	return args.String(0), args.Error(1)
}
