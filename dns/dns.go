package dns

// DnsUpdater can update a DNS record
type DnsUpdater interface {

	// Update the associated DNS record to the given content,
	// returns an error if there was a problem while updating.
	Update(content string) error
}
