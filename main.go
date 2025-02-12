package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the configuration for the application.
type Config struct {
	HotShotURL      string        `json:"hotshot_url"`
	Namespace       uint64        `json:"namespace"`
	PollingInterval time.Duration `json:"polling_interval"`
}

// loadConfig reads the JSON configuration from config/config.json and decodes it into a Config struct.
func loadConfig() Config {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	cfg := loadConfig()
	fmt.Printf("Loaded Config: %+v\n", cfg)

	ticker := time.NewTicker(cfg.PollingInterval * time.Second / 2)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(cfg.HotShotURL + "/status/block-height")
		if err != nil {
			fmt.Println("Error querying block height:", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading block height response:", err)
			continue
		}

		blockHeight, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
		if err != nil {
			fmt.Println("Error parsing block height:", err)
			continue
		}
		fmt.Println("Latest Block Height:", blockHeight)

		availURL := fmt.Sprintf("%s/availability/block/%d/namespace/%d", cfg.HotShotURL, blockHeight, cfg.Namespace)

		availResp, err := http.Get(availURL)
		if err != nil {
			fmt.Println("Error querying availability endpoint:", err)
			continue
		}

		availBody, err := io.ReadAll(availResp.Body)
		availResp.Body.Close()
		if err != nil {
			fmt.Println("Error reading availability response:", err)
			continue
		}
		fmt.Println("Hotshot Namespace Transactions:", string(availBody))
	}
}
