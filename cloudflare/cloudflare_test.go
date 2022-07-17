package cloudflare

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudFlare(t *testing.T) {

	token := "token"
	cf := NewCloudFlare(token)

	assert.NotNil(t, cf.client)

	assert.IsType(t, &httpClient{}, cf.client)

	client := cf.client.(*httpClient)
	assert.Equal(t, token, client.authToken)
}

func TestUpdate(t *testing.T) {

	t.Run("update content", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z", "r")

		cf.client = &mockClient{
			getZone: func(name string) (*Zone, error) {
				return &Zone{Id: "id"}, nil
			},
			getRecord: func(zid, r string) (*Record, error) {
				assert.Equal(t, "id", zid)
				return &Record{Id: "rid", Name: r}, nil
			},
			updateRecord: func(zid string, r *Record) (*Record, error) {
				assert.Equal(t, "id", zid)
				assert.Equal(t, "rid", r.Id)
				assert.Equal(t, "r", r.Name)
				assert.Equal(t, "content", r.Content)
				return &Record{}, nil
			},
		}

		errs := cf.Update("content")
		assert.Len(t, errs, 0)

	})

	t.Run("get zone error", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z1", "r1")
		cf.AddEntry("z1", "r2")
		cf.AddEntry("z2", "r1")

		cf.client = &mockClient{
			getZone: func(name string) (*Zone, error) {
				return nil, errors.New("test error")
			},
		}

		errs := cf.Update("")

		assert.Len(t, errs, 2)
		assert.ErrorContains(t, errs[0], "test error")
		assert.ErrorContains(t, errs[1], "test error")
	})

	t.Run("zone not found", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z1", "r1")
		cf.AddEntry("z2", "r2")

		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {
				return nil, nil
			},
		}

		errs := cf.Update("")

		assert.Len(t, errs, 2)
		assert.ErrorContains(t, errs[0], "not found")
		assert.ErrorContains(t, errs[1], "not found")
	})

	t.Run("get record error", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z1", "r1")
		cf.AddEntry("z2", "r1")

		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {
				return &Zone{}, nil
			},
			getRecord: func(s1, s2 string) (*Record, error) {
				return nil, errors.New("test error")
			},
		}

		errs := cf.Update("")

		assert.Len(t, errs, 2)
		assert.ErrorContains(t, errs[0], "test error")
		assert.ErrorContains(t, errs[1], "test error")

	})

	t.Run("record not found", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z1", "r1")
		cf.AddEntry("z2", "r1")

		cf.client = &mockClient{
			getZone: func(s string) (*Zone, error) {
				return &Zone{}, nil
			},
			getRecord: func(s1, s2 string) (*Record, error) {
				return nil, nil
			},
		}

		errs := cf.Update("")

		assert.Len(t, errs, 2)
		assert.ErrorContains(t, errs[0], "not found")
		assert.ErrorContains(t, errs[1], "not found")

	})

	t.Run("update record error", func(t *testing.T) {

		cf := NewCloudFlare("")
		cf.AddEntry("z1", "r1")
		cf.AddEntry("z2", "r1")

		cf.client = &mockClient{
			getZone: func(name string) (*Zone, error) {
				return &Zone{}, nil
			},
			getRecord: func(zid, r string) (*Record, error) {
				return &Record{}, nil
			},
			updateRecord: func(zid string, r *Record) (*Record, error) {
				return nil, errors.New("test error")
			},
		}

		errs := cf.Update("content")

		assert.Len(t, errs, 2)
		assert.ErrorContains(t, errs[0], "test error")
		assert.ErrorContains(t, errs[1], "test error")

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
