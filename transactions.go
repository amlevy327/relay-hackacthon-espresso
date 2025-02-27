package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func searchLatestTransactions(cfg Config, client *ethclient.Client, lastBlockNumber *uint64) {
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

	fmt.Printf("Searching for transaction at block number %d\n", blockNumber)
	for i, tx := range block.Transactions() {
		if i == 0 {
			continue
		}
		inspectTransaction(tx, cfg)
	}

	*lastBlockNumber = blockNumber
}

func inspectTransaction(tx *types.Transaction, cfg Config) {
	msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Printf("Failed to get sender for transaction %s: %v", tx.Hash().Hex(), err)
		return
	}

	if tx.Value().Int64() >= int64(cfg.Value) && msg.Hex() == cfg.From {
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
}
