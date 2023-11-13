package ip

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	IpifyUrl              = "https://api.ipify.org/"
	QueryParamFormat      = "format"
	QueryParamFormatValue = "text"
	ErrInvalidIp          = "response is not a valid ip: %s"
)

type Ipify struct{}

func (i *Ipify) GetPublicIp() (string, error) {
	r := resty.New().
		SetTransport(&http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
		}).
		SetRetryCount(3)
	res, err := r.NewRequest().
		SetQueryParam(QueryParamFormat, QueryParamFormatValue).
		Get(IpifyUrl)
	if err != nil {
		return "", err
	}
	ip := res.String()
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf(ErrInvalidIp, err)
	}
	return ip, nil
}
