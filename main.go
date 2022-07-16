package main

import (
	"log"
	"time"

	"github.com/garricasaurus/ipflare/ipdetect"
)

func main() {

	// zone := "okki.hu"
	// record := "okki.hu"

	// var updater dns.DnsUpdater = cloudflare.NewCloudFlare(zone, record)

	// err := updater.Update("1.1.1.1")
	// if err != nil {
	// 	panic(err)
	// }

	cd := ipdetect.NewChangeDetector(5 * time.Second)
	cd.Start()

	for {
		select {
		case err := <-cd.Err:
			log.Println(err)
		case ip := <-cd.C:
			log.Printf("ip change: %s", ip)
		}
	}

}
