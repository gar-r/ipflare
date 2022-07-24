package main

import (
	"log"
	"strings"
	"time"

	"ipflare/cloudflare"
	"ipflare/ipdetect"
)

func main() {
	initArgs()
	cd := ipdetect.NewChangeDetector(time.Duration(freq) * time.Second)
	cf := cloudflare.NewCloudFlare(authToken)
	addEntries(cf)
	logStartup()
	startLoop(cd, cf)
}

func startLoop(cd *ipdetect.ChangeDetector, cf *cloudflare.CloudFlare) {
	cd.Start()
	for {
		select {
		case err := <-cd.Err:
			log.Println(err)
		case ip := <-cd.C:
			log.Printf("ip change detected: %s", ip)
			errs := cf.Update(ip)
			if len(errs) > 0 {
				logErrs(errs)
			}
		}
	}
}

func addEntries(cf *cloudflare.CloudFlare) {
	for _, e := range entries {
		parts := strings.Split(e, "/")
		if len(parts) != 2 {
			log.Fatalf("invalid entry: %s", e)
		}
		cf.AddEntry(parts[0], parts[1])
	}
}

func logStartup() {
	log.Println("ipflare starting with the following parameters:")
	log.Printf("%10s: %s", "auth token", "[...]")
	log.Printf("%10s: %d", "frequency", freq)
	log.Printf("%10s: %s", "entries", entries.String())
}

func logErrs(errs []error) {
	for _, err := range errs {
		log.Printf("[error]: %s", err)
	}
}
