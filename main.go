package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 3 BTC to Ivan")
	bc.AddBlock("Send 4 more BTC to Ivan")

	iterator := bc.Iterator()
	for {
		if IsByteEmpty(iterator.currentHash) {
			break
		}

		block := iterator.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n",
			strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}