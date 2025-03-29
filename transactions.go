package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/exec"

	"hackathon-example/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func searchLatestTransactions(cfg config.Config, client *ethclient.Client, lastBlockNumber *uint64) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Failed to get the latest block number: %v", err)
		return
	}

	if blockNumber == *lastBlockNumber {
		return
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Printf("Failed to get block: %v", err)
		return
	}

	log.Printf("Searching for transaction at last processed block number: %d\n", blockNumber)
	for i, tx := range block.Transactions() {
		// Transaction indexes start at 1. Tx 0 is an empty transaction.
		if i == 0 {
			continue
		}
		inspectTransaction(tx, cfg)
	}

	*lastBlockNumber = blockNumber
}

func inspectTransaction(tx *types.Transaction, cfg config.Config) {
	msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Printf("Failed to get sender for transaction %s: %v", tx.Hash().Hex(), err)
		return
	}
	
	if tx.Value().Int64() >= int64(cfg.Value) && tx.To().Hex() == cfg.To {
		log.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
		log.Printf("  Value: %d\n", tx.Value().Int64())
		log.Printf("  To: %s\n", tx.To().Hex())
		log.Printf("  From: %s\n", msg.Hex())
		
		if tx.To() != nil {
			log.Printf("  To: %s\n", tx.To().Hex())
		} else {
			log.Printf("  To: Contract Creation")
		}
		log.Println("---------------------------------------------------------------------------------")

		log.Println("Minting NFT for eligible transaction...")

		// Verify working directory
		wd, _ := os.Getwd()
		log.Printf("Current Working Directory: %s\n", wd)

		// Debug if the script exists
		scriptPath := "./scripts/mintNft.js"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			log.Printf("ERROR: Script %s not found!\n", scriptPath)
			return
		}

		// Run the script
		cmd := exec.Command("npx", "hardhat", "run", "--network", "sepolia", scriptPath)
		cmd.Env = append(os.Environ(), "MINT_TO="+msg.Hex())

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to execute mintNft.js: %v\n", err)
			log.Printf("Script error output: %s\n", string(output))
			return
		}

		log.Printf("mintNft.js Output: %s\n", string(output))
		log.Println("---------------------------------------------------------------------------------")
	}
}
