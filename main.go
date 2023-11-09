package main

import (
	"ipflare/config"
	"ipflare/dns"
	"ipflare/ip"
	"log"
	"os"
	"time"
)

const (
	EnvConfigPath = "CONFIG_PATH"
	EnvApiKey     = "CLOUDFLARE_API_TOKEN"
)

// test
const DefaultConfigPath = "/etc/ipflare/config.yaml"

func main() {
	cfg := initConfig()
	cd := ip.NewChangeDetector(time.Second * time.Duration(cfg.Frequency))
	updater, err := dns.NewCloudflareUpdater(cfg)
	if err != nil {
		panic(err)
	}
	startLoop(cd, updater)
}

func startLoop(cd *ip.ChangeDetector, updater dns.Updater) {
	cd.Start()
	for {
		select {
		case err := <-cd.Err:
			log.Println(err)
		case ip := <-cd.C:
			log.Printf("ip change detected: %s", ip)
			errs := updater.Update(ip)
			if len(errs) > 0 {
				logErrs(errs)
			}
		}
	}
}

func initConfig() *config.Config {
	path := os.Getenv(EnvConfigPath)
	if path == "" {
		path = DefaultConfigPath
	}
	cfg, err := config.Parse(path)
	if err != nil {
		panic(err)
	}
	apiKey := os.Getenv(EnvApiKey)
	if apiKey != "" {
		cfg.ApiToken = apiKey
	}
	logStartup(path, cfg)
	return cfg
}

func logStartup(configLocation string, cfg *config.Config) {
	log.Println("ipflare starting with the following parameters:")
	log.Printf("config location: %s\n", configLocation)
	log.Printf("configuration: %v\n", cfg)
}

func logErrs(errs []error) {
	for _, err := range errs {
		log.Println(err)
	}
}
