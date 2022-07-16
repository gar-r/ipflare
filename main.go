package main

import (
	"flag"
	"log"
	"time"

	"github.com/garricasaurus/ipflare/cloudflare"
	"github.com/garricasaurus/ipflare/ipdetect"
)

var freq int
var zone, record, authToken string

func main() {

	initArgs()

	cf := cloudflare.NewCloudFlare(zone, record, authToken)
	cd := ipdetect.NewChangeDetector(time.Duration(freq) * time.Second)

	logStartup()

	cd.Start()

	for {
		select {
		case err := <-cd.Err:
			log.Println(err)
		case ip := <-cd.C:
			log.Printf("ip change detected: %s", ip)
			err := cf.Update(ip)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func initArgs() {
	flag.IntVar(&freq, "f", 30, "ip change detection frequency in seconds")
	flag.StringVar(&zone, "z", "", "cloudflare zone name")
	flag.StringVar(&record, "r", "", "cloudflare record name")
	flag.StringVar(&authToken, "t", "", "cloudflare api auth token")
	flag.Parse()
}

func logStartup() {
	log.Println("ipflare starting with the following parameters:")
	log.Printf("%10s: %ds", "frequency", freq)
	log.Printf("%10s: %s", "zone", zone)
	log.Printf("%10s: %s", "record", record)
	log.Printf("%10s: %s", "auth token", "[...]")
}
