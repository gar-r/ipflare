package cloudflare

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudFlare(t *testing.T) {

	zone := "zone"
	record := "record"

	t.Run("struct initialized", func(t *testing.T) {
		cf := NewCloudFlare(zone, record)
		assert.Equal(t, zone, cf.zone)
		assert.Equal(t, record, cf.record)
		assert.NotNil(t, cf.client)
	})

	t.Run("auth token obtained from env", func(t *testing.T) {
		token := "test token"
		os.Setenv(authTokenEnvVarName, token)
		cf := NewCloudFlare(zone, record)

		assert.IsType(t, &httpClient{}, cf.client)

		client := cf.client.(*httpClient)
		assert.Equal(t, token, client.authToken)
	})

}

func TestUpdate(t *testing.T) {

	cf := &CloudFlare{
		zone:   "zone",
		record: "record",
	}

	t.Run("get zone error", func(t *testing.T) {
		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {
				return nil, errors.New("test error")
			},
		}

		err := cf.Update("")

		assert.ErrorContains(t, err, "test error")
	})

	t.Run("get record error", func(t *testing.T) {
		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {
				return &Zone{}, nil
			},
			getRecord: func(s1, s2 string) (*Record, error) {
				return nil, errors.New("test error")
			},
		}

		err := cf.Update("")

		assert.ErrorContains(t, err, "test error")
	})

	t.Run("update called with content", func(t *testing.T) {

		content := "content"

		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {

				return &Zone{}, nil
			},
			getRecord: func(s1, s2 string) (*Record, error) {
				return &Record{}, nil
			},
			updateRecord: func(s string, r *Record) (*Record, error) {
				assert.Equal(t, content, r.Content)
				return &Record{}, nil
			},
		}

		err := cf.Update("content")

		assert.NoError(t, err)
	})

}

type mockClient struct {
	getZone      func(string) (*Zone, error)
	getRecord    func(string, string) (*Record, error)
	updateRecord func(string, *Record) (*Record, error)
}

func (m *mockClient) GetZone(zone string) (*Zone, error) {
	return m.getZone(zone)
}

func (m *mockClient) GetRecord(zoneId, recordName string) (*Record, error) {
	return m.getRecord(zoneId, recordName)
}

func (m *mockClient) UpdateRecord(zoneId string, record *Record) (*Record, error) {
	return m.updateRecord(zoneId, record)
}
