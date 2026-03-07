package main

import (
	"encoding/base64"
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

func CreateUULEFromString(location string) (string, error) {
	uuleString, err := encodeUULEString(location)
	if err != nil {
		return "", err
	}
	return uuleString, nil
}

var lengthPrefixKey = "EFGHIjKLMN0PQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789- ABCDEFGHIJKLMNOPQRSTUVWXYZL"

func getLengthPrefix(length int) (string, bool) {
	if length < 4 || length > 89 {
		return "", false
	}
	return string(lengthPrefixKey[length-4]), true
}

func encodeUULEString(uuleString string) (string, error) {
	encodedLocation := base64.StdEncoding.EncodeToString([]byte(uuleString))
	lengthPrefix, exists := getLengthPrefix(len(uuleString))
	if !exists {
		return "", fmt.Errorf("no length prefix found for length %d", len(uuleString))
	}
	return "w+CAIQICI" + lengthPrefix + encodedLocation, nil
}

func scrape(baseURL, apiKey, googleURL, targetOS string) (string, error) {
	params := url.Values{}
	params.Set("apikey", apiKey)
	params.Set("url", googleURL)
	params.Set("js_render", "true")
	params.Set("target_os", targetOS)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result ScrapeResult
	json.Unmarshal(body, &result)

	if result.Result.Status == "completed" {
		return result.Result.Result, nil
	}
	return "", fmt.Errorf("%s", result.Result.ErrorMessage)
}

func main() {
	apiKey := "<YOUR_API_KEY>"
	baseURL := "http://localhost:8082/v1/"

	queries := []struct {
		Query, Region, Location string
	}{
		{"vpn", "us", "Chicago, IL"},
		{"insurance", "us", "New York, NY"},
		{"best credit cards", "gb", "London"},
		{"cheapest flights to tokyo", "us", "Austin, TX"},
		{"hotels in paris", "fr", "Paris"},
		{"restaurants near me", "br", "Sao Paulo"},
		{"best ramen", "jp", "Tokyo"},
	}

	devices := []struct {
		TargetOS string
		Label    string
	}{
		{"android", "mobile"},
		{"linux", "desktop"},
	}

	for _, q := range queries {
		uuleParam, _ := CreateUULEFromString(q.Location)
		googleURL := fmt.Sprintf(
			"https://www.google.com/search?q=%s&gl=%s&uule=%s",
			url.QueryEscape(q.Query),
			q.Region,
			uuleParam,
		)

		fmt.Printf("[%s] => [%s]\n", q.Query, googleURL)

		for _, d := range devices {
			html, err := scrape(baseURL, apiKey, googleURL, d.TargetOS)
			if err != nil {
				fmt.Printf("  [FAIL] %s (%s): %v\n", q.Query, d.Label, err)
				continue
			}

			filename := fmt.Sprintf(
				"serp_%s_%s_%s.html",
				strings.ReplaceAll(q.Query, " ", "_"),
				q.Region,
				d.Label,
			)
			os.WriteFile(filename, []byte(html), 0644)
			fmt.Printf("  [OK] %s (%s, %s, %s) — %d chars => %s\n",
				q.Query, q.Location, q.Region, d.Label, len(html), filename)
		}
	}
}
