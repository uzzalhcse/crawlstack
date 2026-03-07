package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ScrapeResult struct {
	Result struct {
		Status       string `json:"status"`
		Result       string `json:"result"`
		ErrorMessage string `json:"error_message"`
	} `json:"result"`
}

func main() {
	apiKey := "<YOUR_API_KEY>"
	baseURL := "http://localhost:8082/v1/"

	queries := []struct {
		Query, Region, Location string
	}{
		{"vpn", "us", "Chicago, IL"},
		{"insurance", "us", "New York, NY"},
	}

	for _, q := range queries {
		googleURL := fmt.Sprintf(
			"https://www.google.com/search?q=%s&gl=%s&uule=w+CAIQICI%s",
			url.QueryEscape(q.Query), q.Region, url.QueryEscape(q.Location),
		)

		params := url.Values{}
		params.Set("apikey", apiKey)
		params.Set("url", googleURL)
		params.Set("js_render", "true")

		params.Set("target_os", "android")

		resp, err := http.Get(baseURL + "?" + params.Encode())
		if err != nil {
			fmt.Printf("[ERR] %s: %v\n", q.Query, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result ScrapeResult
		json.Unmarshal(body, &result)

		if result.Result.Status == "completed" {
			fmt.Printf("[OK] %s (%s) — %d chars\n", q.Query, q.Location, len(result.Result.Result))
			filename := fmt.Sprintf("serp_%s_%s.html", strings.ReplaceAll(q.Query, " ", "_"), q.Region)
			os.WriteFile(filename, []byte(result.Result.Result), 0644)
		} else {
			fmt.Printf("[FAIL] %s: %s\n", q.Query, result.Result.ErrorMessage)
		}
	}
}
