package dns

import (
	"fmt"
	"ipflare/config"
	"log"
	"time"
)

type Updater interface {
	Update(ip string) []error
}

type CloudflareUpdater struct {
	Client   Client
	EntrySet config.EntrySet
}

func NewCloudflareUpdater(cfg *config.Config) (Updater, error) {
	client, err := NewCloudflareClient(cfg.ApiToken)
	if err != nil {
		return nil, err
	}
	return &CloudflareUpdater{
		Client:   client,
		EntrySet: cfg.Entries,
	}, nil
}

func (c *CloudflareUpdater) Update(ip string) []error {
	errors := make([]error, 0)
	for zoneName := range c.EntrySet {
		records, err := c.Client.GetDNSRecords(zoneName)
		if err != nil {
			errors = append(errors, err)
			return errors
		}
		for _, record := range records {
			if c.isManaged(zoneName, record.Name) {
				record.Content = ip
				record.Comment = c.makeComment()
				log.Printf("updating %s", record.Name)
				err := c.Client.UpdateDNSRecord(record)
				if err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	return errors
}

func (c *CloudflareUpdater) isManaged(zoneName, recordName string) bool {
	for _, entry := range c.EntrySet[zoneName] {
		if entry == recordName {
			return true
		}
	}
	return false
}

func (c *CloudflareUpdater) makeComment() string {
	return fmt.Sprintf("last updated by ipflare at %s", time.Now().String())
}
