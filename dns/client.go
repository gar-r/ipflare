package dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

type Client interface {
	GetDNSRecords(zoneName string) ([]*Record, error)
	UpdateDNSRecord(record *Record) error
}

type Record struct {
	ZoneId  string
	ID      string
	Name    string
	Type    string
	Content string
	Comment string
	Proxied bool
}

type CloudflareClient struct {
	api *cloudflare.API
}

func NewCloudflareClient(apiToken string) (Client, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	return &CloudflareClient{api}, err
}

func (c *CloudflareClient) GetDNSRecords(zoneName string) ([]*Record, error) {
	zoneId, err := c.api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, err
	}
	records, _, err := c.api.ListDNSRecords(context.Background(),
		cloudflare.ZoneIdentifier(zoneId),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
		})
	if err != nil {
		return nil, err
	}
	result := make([]*Record, len(records))
	for i, record := range records {
		result[i] = &Record{
			ZoneId:  zoneId,
			ID:      record.ID,
			Name:    record.Name,
			Type:    record.Type,
			Content: record.Content,
			Proxied: *record.Proxied,
			Comment: record.Comment,
		}
	}
	return result, nil
}

func (c *CloudflareClient) UpdateDNSRecord(record *Record) error {
	_, err := c.api.UpdateDNSRecord(context.Background(),
		cloudflare.ZoneIdentifier(record.ZoneId),
		cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Name:    record.Name,
			Type:    record.Type,
			Content: record.Content,
			Proxied: cloudflare.BoolPtr(record.Proxied),
			Comment: cloudflare.StringPtr(record.Comment),
		})
	return err
}
