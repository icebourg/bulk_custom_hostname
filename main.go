package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/exp/slices"
)

// ugh, I wish the Cloudflare SDK didn't make me implement pagination
func ListCustomHostnames(ctx context.Context, api *cloudflare.API, zoneId string) ([]string, *cloudflare.ResultInfo, error) {
	var hostnames []string
	var lastResultInfo cloudflare.ResultInfo

	page := 0

	for {
		customHostnames, results, err := api.CustomHostnames(ctx, zoneId, page, cloudflare.CustomHostname{})
		if err != nil {
			log.Fatal(err)
		}

		for _, h := range customHostnames {
			hostnames = append(hostnames, h.Hostname)
		}
		lastResultInfo = results

		page = results.Next().Page
		if results.Done() {
			break
		}
	}

	return hostnames, &lastResultInfo, nil
}

func main() {
	batchSize := flag.Int("batchsize", 100, "the number of domain names we will operate on at a given time")
	zoneId := flag.String("zoneid", "REQUIRED", "the zone id of the SSL/SaaS application")
	exclusions := flag.String("exclusions", "", "A comma separated list of domains to exclude from consideration (eg example.com,example2.com)")
	flag.Parse()

	if *zoneId == "REQUIRED" {
		log.Fatal("Please pass required zoneID value")
	}

	excludedDomains := strings.Split(*exclusions, ",")

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch Zones
	zones, err := api.ListZones(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// fetch existing custom SSL/TLS hostnames
	customHostnames, _, err := ListCustomHostnames(ctx, api, *zoneId)
	if err != nil {
		log.Fatal(err)
	}

	var domainRequests []string

	for _, z := range zones {
		if slices.Index[string](excludedDomains, z.Name) > -1 {
			// excluding this domain because we were asked to exclude it
			continue
		}

		if slices.Index[string](customHostnames, z.Name) > -1 {
			// excluding this domain, we already have a certificate issued for it
			continue
		}

		fmt.Printf("Custom hostname needed for %s \n", z.Name)
		domainRequests = append(domainRequests, z.Name)

		if len(domainRequests) >= *batchSize {
			requestHostnames(ctx, api, domainRequests, *zoneId)
			domainRequests = nil
		}
	}

	// collect all remaining hostnames
	if len(domainRequests) > 0 {
		requestHostnames(ctx, api, domainRequests, *zoneId)
		domainRequests = nil
	}
}

func requestHostnames(ctx context.Context, api *cloudflare.API, domainRequests []string, zoneId string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Ready to request %d custom hostnames? Yes to proceed, any other value to quit...", len(domainRequests))
	s, _ := r.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(s)) != "yes" {
		os.Exit(5)
	}

	for _, hostname := range domainRequests {
		fmt.Printf("Requesting %s\n", hostname)
		sslOptions := cloudflare.CustomHostnameSSL{
			Method: "txt",
			Type:   "dv",
		}

		_, err := api.CreateCustomHostname(ctx, zoneId, cloudflare.CustomHostname{Hostname: hostname, SSL: &sslOptions})
		if err != nil {
			log.Fatal(err)
		}
	}
}
