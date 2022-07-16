package cloudflare

import "os"

const authTokenEnvVarName = "CLOUDFLARE_API_TOKEN"

// CloudFlare implements Updater and uses the CloudFlare API to modify
// a DNS record identified by its zone name and record name.
type CloudFlare struct {
	client client
	zone   string
	record string
}

func NewCloudFlare(zone, record string) *CloudFlare {
	return &CloudFlare{
		zone:   zone,
		record: record,
		client: &httpClient{
			authToken: os.Getenv(authTokenEnvVarName),
		},
	}
}

// Update the content of the associated record.
func (u *CloudFlare) Update(content string) error {
	z, err := u.client.GetZone(u.zone)
	if err != nil {
		return err
	}
	r, err := u.client.GetRecord(z.Id, u.record)
	if err != nil {
		return err
	}
	r.Content = content
	_, err = u.client.UpdateRecord(z.Id, r)
	return err
}
