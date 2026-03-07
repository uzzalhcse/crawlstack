# CrawlStack

A self-hosted universal web scraper API. Extract data from any website with a single `GET /v1/` request — returns clean HTML, Markdown, or screenshots.

## Features

- **HTTP Fast Path** — TLS-spoofed requests without launching a browser
- **Browser Engine** — Real Chrome for JS-heavy & protected sites
- **Stealth+** — Fingerprint spoofing & virtual display bypass
- **JS Instructions** — Click, fill, scroll, wait via JSON directives
- **Session Persistence** — Reuse browser sessions across requests
- **Proxy Routing** — Built-in premium proxy support
- **CAPTCHA Solving** — Auto-solve with `solve_captcha=true`
- **Multiple Output Formats** — HTML, Markdown, or screenshot

## Quick Start

```bash
# 1. Start all services
docker compose -f docker-compose.prod.yml up -d

# 2. Seed the database (first time only)
docker compose -f docker-compose.prod.yml exec backend ./app seed

# 3. Open the dashboard at http://localhost:3232
#    Default credentials: admin@crawlstack.com / password
```

## API Usage

```
GET /v1/?apikey=<YOUR_API_KEY>&url=<TARGET_URL>
```

### Parameters

| Parameter | Description |
|---|---|
| `apikey` | Your API key (required) |
| `url` | Target URL to scrape (required) |
| `js_render` | Enable headless Chrome rendering (`true`/`false`) |
| `output_format` | `html` (default), `markdown`, or `screenshot` |
| `wait_for` | CSS selector to wait for before returning |
| `wait_time` | Extra wait time in ms after selector match |
| `js_instructions` | JSON array of browser actions (click, fill, wait, evaluate) |
| `session_id` | Reuse a browser session across requests |
| `premium_proxy` | Route through premium proxy (`true`/`false`) |
| `solve_captcha` | Auto-solve CAPTCHAs (`true`/`false`) |
| `target_os` | Emulate device fingerprint (`windows`, `android`, etc.) |

### Example

```bash
# Basic scrape
curl "http://localhost:8082/v1/?apikey=<YOUR_API_KEY>&url=https%3A%2F%2Fexample.com"

# JS rendering with wait
curl "http://localhost:8082/v1/?apikey=<YOUR_API_KEY>&url=https%3A%2F%2Fexample.com&js_render=true&wait_for=.content"
```

## Architecture

| Service | Image |
|---|---|
| **backend** | `ghcr.io/camoium/crawlstack/backend` |
| **scraper** | `ghcr.io/camoium/crawlstack/scraper` |
| **frontend** | `ghcr.io/camoium/crawlstack/frontend` |
| **postgres** | `postgres:16-alpine` |
| **redis** | `redis:7-alpine` |

## Environment Variables

```env
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=crawlstack
JWT_SECRET=your-random-secret-key
JWT_EXPIRY=72h
FRONTEND_URL=http://localhost:3232
REDIS_PASSWORD=
TZ=America/New_York
WORKER_COUNT=4
WARM_POOL_SIZE=2
BROWSER_TIMEOUT=120s
```
