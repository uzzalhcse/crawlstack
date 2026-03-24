# Google SERP Scraper

A Go script that queries Google Search across multiple locations using UULE encoding and saves the HTML results locally.

## Prerequisites

- A running scraper API at `http://localhost:8083/scrape/`
- A valid API key for that service
- **Docker** (recommended), or [Go](https://golang.org/dl/) 1.18+ for running directly

## Setup

Open `main.go` and replace the placeholder with your API key:
```go
apiKey := "<YOUR_API_KEY>"
```

If your scraper API runs on a different address, update `baseURL` accordingly.

---

## Run with Docker (recommended)

**1. Build the image:**
```bash
docker build -t serp-scraper .
```

**2. Run the container:**
```bash
docker run --rm \
  --network host \
  -v $(pwd)/output:/output \
  serp-scraper
```

- `--network host` lets the container reach your locally running scraper API on `localhost:8083`
- `-v $(pwd)/output:/output` mounts a local `output/` folder so the generated HTML files are saved on your machine

> **Windows (PowerShell):** Replace `$(pwd)` with `${PWD}`.

---

## Run without Docker

```bash
go run main.go
```

---
## Example Urls to use in Crawlstack

```
[vpn]=>[https://www.google.com/search?q=vpn&gl=us&uule=w+CAIQICILQ2hpY2FnbywgSUw=]
[insurance]=>[https://www.google.com/search?q=insurance&gl=us&uule=w+CAIQICIMTmV3IFlvcmssIE5Z]
[best credit cards]=>[https://www.google.com/search?q=best+credit+cards&gl=gb&uule=w+CAIQICIGTG9uZG9u]
[cheapest flights to tokyo]=>[https://www.google.com/search?q=cheapest+flights+to+tokyo&gl=us&uule=w+CAIQICIKQXVzdGluLCBUWA==]
[hotels in paris]=>[https://www.google.com/search?q=hotels+in+paris&gl=fr&uule=w+CAIQICIFUGFyaXM=]
[restaurants near me]=>[https://www.google.com/search?q=restaurants+near+me&gl=br&uule=w+CAIQICIjU2FvIFBhdWxv]
[best ramen]=>[https://www.google.com/search?q=best+ramen&gl=jp&uule=w+CAIQICIFVG9reW8=]
```

## What It Does

For each query defined in the `queries` slice, the script runs **two requests** — one for mobile (`android`) and one for desktop (`linux`):

1. Encodes the target location as a [UULE parameter](https://moz.com/ugc/geolocation-the-hidden-localization-factor)
2. Builds a Google Search URL with the query, region (`gl`), and location (`uule`)
3. Sends the URL to the scraper API twice — once per device type — with JS rendering enabled
4. Saves each result as a separate HTML file in the working directory

**Output files** are named `serp_<query>_<region>_<device>.html`, e.g.:
```
serp_vpn_us_mobile.html
serp_vpn_us_desktop.html
serp_best_credit_cards_gb_mobile.html
serp_best_credit_cards_gb_desktop.html
```

---

## Adding Queries

Extend the `queries` slice in `main.go`:

```go
{"your query", "country_code", "City, State"},
```

Country codes follow the [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) standard (e.g. `us`, `gb`, `fr`).

> After editing `main.go`, rebuild the Docker image before running again.