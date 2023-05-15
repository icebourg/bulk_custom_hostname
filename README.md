# bulk_custom_hostname

This little script will help you bulk request [custom hostnames](https://developers.cloudflare.com/cloudflare-for-platforms/cloudflare-for-saas/) from Cloudflare.

## Usage
Create a [*scoped* API key](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) within Cloudflare with two permissions:

1. Permission to manage/read your DNS zones
2. Permission to manage your SSL/TLS settings.

We will attempt to create custom hostnames for *every* DNS zone we find we have access to, so to limit that to a particular account or particular zones, use API scopes to limit what this script has access to.

Your scoped API key should be available under the environment variable `CLOUDFLARE_API_KEY`. 

### Arguments
- `batchsize` The number of domains we will request at a given time
- `zoneid` (REQUIRED) The ID of the Zone you want to create the records under.
- `exclusions` A comma separated list of domains to exclude from your requests

### Example Usage
```
bulk_custom_hostname --zoneid EXAMPLE12345 --exclusions example.com,exampledev.com --batchsize 10
./ --zoneid 821a202c922d70a38c1eaaca225f3007  --exclusions 22onerealty.com,4stgeorgerealestate.com
Custom hostname needed for example1.com
Custom hostname needed for example2.com
Custom hostname needed for example3.com
Custom hostname needed for example4.com
Custom hostname needed for example5.com
Custom hostname needed for example6.com
Custom hostname needed for example7.com
Custom hostname needed for example8.com
Custom hostname needed for example9.com
Custom hostname needed for example10.com
Ready to request 10 custom hostnames? Yes to proceed, any other value to quit...yes
Requesting example1.com
Requesting example2.com
Requesting example3.com
Requesting example4.com
Requesting example5.com
Requesting example6.com
Requesting example7.com
Requesting example8.com
Requesting example9.com
Requesting example10.com
...
```

### Installation
Download and install [Go](https://go.dev/doc/install), then it's as easy as:

```
% go install github.com/icebourg/bulk_custom_hostname@latest
go: downloading github.com/icebourg/bulk_custom_hostname 
% bulk_custom_hostname --help
Usage of bulk_custom_hostname:
  -batchsize int
    	the number of domain names we will operate on at a given time (default 100)
  -exclusions string
    	A comma separated list of domains to exclude from consideration (eg example.com,example2.com)
  -zoneid string
    	the zone id of the SSL/SaaS application (default "REQUIRED")
```

If `bulk_custom_hostname` doesn't find the binary after a successful `go install ...`, make sure your Go installation directory is in your `$PATH`. Run `go env` to find where go is installing your binaries.
