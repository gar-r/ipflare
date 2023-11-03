package dns

import (
	"errors"
	"ipflare/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCloudflareUpdater_Update(t *testing.T) {
	t.Run("error in get dns records", func(t *testing.T) {
		client := &TestClient{}
		client.On("GetDNSRecords", "zone").Return(make([]*Record, 0), errors.New("test error"))
		updater := &CloudflareUpdater{
			Client:   client,
			EntrySet: config.EntrySet{"zone": []string{}},
		}

		errors := updater.Update("127.0.0.1")

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "test error", errors[0].Error())
	})

	t.Run("error in update dns records", func(t *testing.T) {
		client := &TestClient{}
		r1 := &Record{Name: "foo.example.com"}
		r2 := &Record{Name: "bar.example.com"}
		client.On("GetDNSRecords", "zone").Return([]*Record{r1, r2}, nil)
		client.On("UpdateDNSRecord", r1).Return(errors.New("error1"))
		client.On("UpdateDNSRecord", r2).Return(errors.New("error2"))
		updater := &CloudflareUpdater{
			Client: client,
			EntrySet: config.EntrySet{
				"zone": []string{
					"foo.example.com",
					"bar.example.com",
				},
			},
		}

		errors := updater.Update("127.0.0.1")

		assert.Equal(t, 2, len(errors))
		assert.Equal(t, "error1", errors[0].Error())
		assert.Equal(t, "error2", errors[1].Error())
	})

	t.Run("update managed records", func(t *testing.T) {
		client := &TestClient{}
		r1 := &Record{Name: "foo.example.com"}
		r2 := &Record{Name: "bar.example.com"}
		r3 := &Record{Name: "baz.example.com"}
		client.On("GetDNSRecords", "zone").Return([]*Record{r1, r2}, nil)
		client.On("UpdateDNSRecord", r1).Return(nil)
		client.On("UpdateDNSRecord", r2).Return(nil)
		updater := &CloudflareUpdater{
			Client: client,
			EntrySet: config.EntrySet{
				"zone": []string{
					"foo.example.com",
					"bar.example.com",
					"baz.example.com",
				},
			},
		}

		errors := updater.Update("127.0.0.1")

		assert.Equal(t, 0, len(errors))
		client.AssertNotCalled(t, "UpdateDNSRecord", r3)
	})
}

type TestClient struct {
	mock.Mock
}

func (t *TestClient) GetDNSRecords(zoneName string) ([]*Record, error) {
	args := t.Called(zoneName)
	return args.Get(0).([]*Record), args.Error(1)
}

func (t *TestClient) UpdateDNSRecord(record *Record) error {
	args := t.Called(record)
	return args.Error(0)
}
