package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	CaffNodeURL     string        `json:"caff_node_url"`
	PollingInterval time.Duration `json:"polling_interval"`
	Value           int           `json:"value"`
	From            string        `json:"from"`
}

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

	client, err := ethclient.Dial(cfg.CaffNodeURL)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var lastBlockNumber uint64

	for range ticker.C {
		searchLatestTransactions(cfg, client, &lastBlockNumber)
	}
}
