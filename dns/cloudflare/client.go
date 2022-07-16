package cloudflare

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://api.cloudflare.com/client/v4"

// client is a CloudFlare api client
type client interface {

	// GetZone tries to get the Zone with the given name.
	// In case of multiple matches, the first one is returned.
	GetZone(zone string) (*Zone, error)

	// GetRecord tries to get a DNS Record under the given zone id
	// with the given record name.
	// In case of multiple matches, the first one is returned.
	GetRecord(zoneId, recordName string) (*Record, error)

	// UpdateRecord tries to update a DNS Record
	UpdateRecord(zoneId string, record *Record) (*Record, error)
}

type httpClient struct {
	authToken string
}

func (c *httpClient) GetZone(zone string) (*Zone, error) {
	url := fmt.Sprintf("%s/zones", baseURL)
	robj := &Response[[]*Zone]{}
	res, err := c.getRequest().
		SetQueryParam("name", zone).
		SetResult(robj).
		Get(url)
	if err != nil {
		return nil, err
	}
	arr, err := unpack(res, robj)
	return first(arr), err
}

func (c *httpClient) GetRecord(zoneId, recordName string) (*Record, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records", baseURL, zoneId)
	robj := &Response[[]*Record]{}
	res, err := c.getRequest().
		SetQueryParam("name", recordName).
		SetResult(robj).
		Get(url)
	if err != nil {
		return nil, err
	}
	arr, err := unpack(res, robj)
	return first(arr), err
}

func (c *httpClient) UpdateRecord(zoneId string, record *Record) (*Record, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", baseURL, zoneId, record.Id)
	robj := &Response[*Record]{}
	res, err := c.getRequest().
		SetBody(record).
		SetResult(robj).
		Put(url)
	if err != nil {
		return nil, err
	}
	return unpack(res, robj)
}

func (c *httpClient) getRequest() *resty.Request {
	r := resty.New()
	return r.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.authToken)
}

func unpack[T any](httpResp *resty.Response, robj *Response[T]) (T, error) {
	var zero T
	if httpResp.IsError() {
		return zero, errors.New(httpResp.Status())
	}
	if !robj.Success {
		return zero, fmt.Errorf("%v", robj.Errors)
	}
	return robj.Result, nil
}

func first[T any](arr []T) T {
	if len(arr) == 0 {
		var zero T
		return zero
	}
	return arr[0]
}
