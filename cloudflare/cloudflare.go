package cloudflare

import (
	"fmt"
)

// CloudFlare uses the CloudFlare public API to modify DNS entries.
type CloudFlare struct {
	client  client
	entries map[string]map[string]bool
}

// NewCloudFlare creates a new CloudFlare instance with the given auth token.
func NewCloudFlare(authToken string) *CloudFlare {
	return &CloudFlare{
		entries: make(map[string]map[string]bool),
		client: &httpClient{
			authToken: authToken,
		},
	}
}

// AddEntry adds a DNS entry (zone/record) to the list of entries to update.
func (c *CloudFlare) AddEntry(zone, record string) {
	if _, ok := c.entries[zone]; !ok {
		c.entries[zone] = map[string]bool{}
	}
	c.entries[zone][record] = true
}

// Update the content of all registered DNS entries.
// Execution will not abort in case of errors, each error is returned.
func (c *CloudFlare) Update(content string) []error {
	errors := make([]error, 0)
	for zoneName := range c.entries {
		errs := c.updateZone(zoneName, content)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	return errors
}

func (c *CloudFlare) updateZone(zoneName string, content string) []error {
	errors := make([]error, 0)
	zone, err := c.client.GetZone(zoneName)
	if err != nil {
		errors = append(errors, err)
		return errors
	}
	if zone == nil {
		errors = append(errors, fmt.Errorf("zone not found: %s", zoneName))
		return errors
	}
	for recordName := range c.entries[zoneName] {
		err := c.updateRecord(zone.Id, recordName, content)
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}
	return errors
}

func (c *CloudFlare) updateRecord(zoneId, recordName, content string) error {
	record, err := c.client.GetRecord(zoneId, recordName)
	if err != nil {
		return err
	}
	if record == nil {
		return fmt.Errorf("record not found: %s", recordName)
	}
	record.Content = content
	_, err = c.client.UpdateRecord(zoneId, record)
	if err != nil {
		return err
	}
	return nil
}
