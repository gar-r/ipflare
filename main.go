package main

import (
	"github.com/garricasaurus/ipflare/dns"
	"github.com/garricasaurus/ipflare/dns/cloudflare"
)

func main() {

	zone := "okki.hu"
	record := "okki.hu"

	var updater dns.DnsUpdater = cloudflare.NewCloudFlare(zone, record)

	err := updater.Update("1.1.1.1")
	if err != nil {
		panic(err)
	}

}
