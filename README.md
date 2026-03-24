
# Introduction

CrawlStack is a self-hosted web scraping toolkit. You give it a URL, it gives you back the page as HTML, Markdown, or a screenshot. Under the hood it runs the [Camoium](https://github.com/camoium) anti-detection browser, which spoofs fingerprints so target sites see a real browser instead of a bot.

There are two main things you can do with it:

- **Scraping API** — send a GET request, get page content back. Handles JS rendering, CAPTCHA solving, proxies, the whole deal.
- **Browser API** — connect your own automation scripts (puppeteer, playwright, rod, chromedp) to a managed browser pool over CDP WebSocket. You write the automation, CrawlStack handles the browser lifecycle and fingerprinting.

Both services are protected by an API key (`API_KEY` environment variable). Every request must include the key via `?api_key=` query parameter or `X-API-Key` header. The services won't start without it.



## Quick Start

```bash
# 1.Copy the example env file and set your API key
cp .env.example .env
# Edit .env and set API_KEY=your-secret-key

# 2. Start all services
docker compose up -d

# 3. Open the dashboard at http://localhost:3232
#    Default credentials: admin@crawlstack.com / password

```

## Architecture

The system is made up of four services. You don't need all of them — the scraper and browser-api work standalone without the backend or frontend.

| Service | Port | What it does |
|---------|------|-------------|
| **Scraper** | 8083 | The actual scraping engine. Takes a URL, returns content. |
| **Browser API** | 9222 | CDP WebSocket proxy for remote browser automation. |
| **Backend** | 8082 | API gateway with auth. Proxies requests to the scraper. Only needed if you want the web UI. |
| **Frontend** | 3232 | Web dashboard for testing scrapes and managing settings. Optional. |

**Standalone mode** (scraper + browser-api only):

```
Your code ---GET /scrape?api_key=...--> Scraper (:8083) ---> HTTP fetch or Browser
Your code ---WebSocket?api_key=...---> Browser API (:9222) ---> Managed browser session
```

**Full stack mode** (with UI):

```
Browser -----> Frontend (:3232) --nginx proxy--> Backend (:8082) ---> Scraper (:8083)
Your code ---> Browser API (:9222) ---> Managed browser session
```

## Scraping API

Hit `/scrape` with a URL and get the page content back.

```
GET http://localhost:8083/scrape?api_key=your-secret-key&url=https://example.com&js_render=true
```

When `js_render` is off (the default), it uses a fast HTTP client with TLS fingerprinting — no browser needed, responses come back in under 2 seconds. When `js_render=true`, it spins up a headless Camoium browser with full fingerprint spoofing.

You can also pass proxies, wait for CSS selectors, execute JS instructions (clicks, form fills, waits), solve CAPTCHAs, and choose between HTML, Markdown, or screenshot output.

## Browser API

If you need more control than the Scraping API gives you, connect directly to a browser instance over CDP:

```
ws://localhost:9222/?api_key=your-secret-key&target_os=windows&fingerprint_mode=random
```

## Local API Documentation

```
http://localhost:3232/docs/introduction
```
### Live API docs with Playground, examples for the Scraping API and Browser API.

```
http://34.85.113.40:3232/docs
```