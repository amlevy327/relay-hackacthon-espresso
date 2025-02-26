package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	HotShotURL      string        `json:"hotshot_url"`
	ChainID         uint64        `json:"chain_id"`
	PollingInterval time.Duration `json:"polling_interval"`
	CaffNodeURL     string        `json:"caff_node_url"`
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
		transactions(cfg, client, &lastBlockNumber)
	}
}

func transactions(cfg Config, client *ethclient.Client, lastBlockNumber *uint64) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Failed to get the latest block number: %v", err)
		return
	}

	if blockNumber == *lastBlockNumber {
		return // Skip if we've already seen this block
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Printf("Failed to get block: %v", err)
		return
	}

	fmt.Printf("Latest Block Number: %d\n", blockNumber)
	fmt.Println("Transactions in the Latest Block:")

	for i, tx := range block.Transactions() {
		if i == 0 {
			continue
		}
		printTransactionDetails(tx)
	}

	*lastBlockNumber = blockNumber
}

func printTransactionDetails(tx *types.Transaction) {
	msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Printf("Failed to get sender for transaction %s: %v", tx.Hash().Hex(), err)
		return
	}

	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  Value: %d\n", tx.Value().Int64())
	fmt.Printf("  From: %s\n", msg.Hex())
	if tx.To() != nil {
		fmt.Printf("  To: %s\n", tx.To().Hex())
	} else {
		fmt.Println("  To: Contract Creation")
	}
	fmt.Println("---------------------------")
}
