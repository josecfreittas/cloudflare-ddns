package main

import (
	"flag"
	"log"
	"os"
	"time"

	internal "cloudflare-ddns/internal"
)

func checkIPv4() (string, error) {
	return internal.HTTPGet("https://checkip.amazonaws.com", nil)
}

func checkIPv6() (string, error) {
	return internal.HTTPGet("https://v6.ident.me", nil)
}

func main() {
	duration := flag.Duration("duration", 0, "update interval (ex. 15s, 1m, 6h); if not specified or set to 0s, run only once and exit")
	updateIPv4 := flag.Bool("ipv4", true, "update A record")
	updateIPv6 := flag.Bool("ipv6", false, "update AAAA record")
	flag.Parse()

	// Environment variables (moved from provider)
	apiToken := os.Getenv("CF_TOKEN")
	if apiToken == "" {
		log.Fatalf("CF_TOKEN env. variable is required")
	}
	zoneID := os.Getenv("CF_ZONE_ID")
	if zoneID == "" {
		log.Fatalf("CF_ZONE_ID env. variable is required")
	}
	host := os.Getenv("CF_HOST")
	if host == "" {
		log.Fatalf("CF_HOST env. variable is required")
	}

	// Track last observed IPs to avoid unnecessary API calls
	var lastIPv4, lastIPv6 string

	if *duration == time.Duration(0) {
		runDDNS(apiToken, zoneID, host, *updateIPv4, *updateIPv6, &lastIPv4, &lastIPv6)
		return
	}

	// Run once immediately, then on each tick
	runDDNS(apiToken, zoneID, host, *updateIPv4, *updateIPv6, &lastIPv4, &lastIPv6)
	ticker := time.NewTicker(*duration)
	defer ticker.Stop()
	for range ticker.C {
		runDDNS(apiToken, zoneID, host, *updateIPv4, *updateIPv6, &lastIPv4, &lastIPv6)
	}
}

func runDDNS(apiToken, zoneID, host string, updateIPv4, updateIPv6 bool, lastIPv4, lastIPv6 *string) {
	updateIfChanged := func(label string, checkIP func() (string, error), lastObserved *string, updateRecord func(string) error) {
		currentIP, err := checkIP()
		if err != nil {
			log.Printf("failed to check %s: %v", label, err)
			return
		}
		if currentIP == *lastObserved {
			log.Printf("%s unchanged (%s)", label, currentIP)
			return
		}
		*lastObserved = currentIP
		if err := updateRecord(currentIP); err != nil {
			log.Printf("failed to update %s record: %v", label, err)
			return
		}
		log.Printf("%s record updated to %s", label, currentIP)
	}

	if updateIPv4 {
		updateIfChanged("IPv4", checkIPv4, lastIPv4, func(ip string) error { return internal.UpdateRecord(apiToken, zoneID, host, ip, "A") })
	}
	if updateIPv6 {
		updateIfChanged("IPv6", checkIPv6, lastIPv6, func(ip string) error { return internal.UpdateRecord(apiToken, zoneID, host, ip, "AAAA") })
	}
}
